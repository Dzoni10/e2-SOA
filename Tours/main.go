package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"tours/database"
	"tours/handler"
	"tours/metrics"
	"tours/repo"
	"tours/service"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var tp *sdktrace.TracerProvider

const serviceName = "tour-service"

func initTracer() (*sdktrace.TracerProvider, error) {
	url := os.Getenv("JAEGER_ENDPOINT")

	if url == "" {
		url = "http://jaeger:14268/api/traces"
		log.Printf("Using default Jaeger endpoint: %s", url)
	} else {
		log.Printf("Using Jaeger endpoint from env: %s", url)
	}

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))

	if err != nil {
		return nil, err
	}

	// Kreiraj resource sa service info
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(time.Second*5)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	log.Println("Tracer initialized successfully")

	return tp, nil
}

func tracingMiddleWate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer(serviceName)

		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path)
		defer span.End()

		// Dodaj span context u request
		r = r.WithContext(ctx)

		// Pozovi sledeći handler
		next.ServeHTTP(w, r)
	})
}

// Middleware za Prometheus metrics
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Custom ResponseWriter da uhvatimo status kod
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(lrw.statusCode)

		// Zabelezi metrics
		metrics.RequestDuration.WithLabelValues(r.URL.Path).Observe(duration)
		metrics.RequestCount.WithLabelValues(r.URL.Path, r.Method, status).Inc()
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func main() {

	metrics.InitMetrics()
	log.Println("Metrics initialized")

	var err error
	tp, err = initTracer()

	if err != nil {
		log.Printf("Failed to initialize tracer: %v", err)
		log.Println("Continuing without tracing...")
	} else {
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				log.Printf("Error shutting down tracer: %v", err)
			}
		}()
	}

	//http.Handle("/metrics", promhttp.Handler())

	database.Init()

	keypointRepo := &repo.KeyPointRepository{Collection: database.KeyPointCollection}
	keypointSrv := &service.KeyPointService{Repo: keypointRepo}
	keypointHandler := &handler.KeyPointHandler{Service: keypointSrv}

	tourRepo := &repo.TourRepository{}
	srv := &service.TourService{Repo: tourRepo, KeyPointRepo: keypointRepo}
	h := &handler.TourHandler{Service: srv}

	reviewRepo := &repo.ReviewRepository{}
	reviewSrv := &service.ReviewService{Repo: reviewRepo}
	reviewHandler := &handler.ReviewHandler{Service: reviewSrv}

	positionRepo := &repo.PositionRepository{}
	positionSrv := &service.PositionService{Repo: positionRepo}
	positionHandler := &handler.PositionHandler{Service: positionSrv}

	executionRepo := &repo.TourExecutionRepository{}
	executionSrv := &service.TourExecutionService{Repo: executionRepo}
	executionHandler := &handler.TourExecutionHandler{Service: executionSrv}

	r := mux.NewRouter()

	r.Use(tracingMiddleWate)
	r.Use(metricsMiddleware)

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/tours/all", h.GetAllTours).Methods("GET")
	r.HandleFunc("/tours/{id}", h.GetTour).Methods("GET")
	r.HandleFunc("/tours", h.CreateTour).Methods("POST")
	r.HandleFunc("/tours/{id}/update-length", h.UpdateLength).Methods("PUT")

	r.HandleFunc("/tours/{id}/reviews", reviewHandler.GetReviewsForTour).Methods("GET")
	r.HandleFunc("/reviews", reviewHandler.CreateReview).Methods("POST")
	r.HandleFunc("/reviews/{id}/add-image", reviewHandler.AddImageToReview).Methods("POST")
	r.HandleFunc("/reviews/{id}/remove-image", reviewHandler.RemoveImageFromReview).Methods("DELETE")
	r.HandleFunc("/uploads/reviews/{filename}", reviewHandler.ServeImage).Methods("GET")
	r.HandleFunc("/tours/calculate-metrics", h.CalculateMetrics).Methods("POST")
	r.HandleFunc("/tours/keypoints", keypointHandler.CreateKeyPoint).Methods("POST")                               // Changed from /tours/keypoints
	r.HandleFunc("/tours/keypoints/bulk-update-tour-id", keypointHandler.BulkUpdateKeyPointsTourId).Methods("PUT") // Added missing route
	r.HandleFunc("/keypoints/{id}/update-tour-id", keypointHandler.UpdateKeyPointTourId).Methods("PUT")            // Optional individual update
	r.HandleFunc("/tours/keypoints/{id}", keypointHandler.UpdateKeypoint).Methods("PUT")                           // Optional individual update
	r.HandleFunc("/tours/{id}/keypoints", keypointHandler.GetKeyPointsForTour).Methods("GET")
	r.HandleFunc("/keypoints/{id}", keypointHandler.DeleteKeyPoint).Methods("DELETE")
	r.HandleFunc("/keypoints/{id}/order", keypointHandler.UpdateKeyPointOrder).Methods("PUT")
	r.HandleFunc("/keypoints/bulk-reorder", keypointHandler.BulkReorderKeyPoints).Methods("PUT")
	r.HandleFunc("/keypoints/{id}/add-image", keypointHandler.AddImageToKeyPoint).Methods("POST")
	r.HandleFunc("/keypoints/{id}/remove-image", keypointHandler.RemoveImageFromKeyPoint).Methods("DELETE")
	r.HandleFunc("/uploads/keypoints/{filename}", keypointHandler.ServeImage).Methods("GET")

	r.HandleFunc("/position/{id}", positionHandler.GetPosition).Methods("GET")
	r.HandleFunc("/position/{id}/create", positionHandler.InitializePosition).Methods("POST")
	r.HandleFunc("/position/{id}/update", positionHandler.UpdatePosition).Methods("PUT")

	r.HandleFunc("/tours/{id}/status", h.UpdateStatus).Methods("PUT")
	r.HandleFunc("/tours/{id}/status-info", h.GetTourStatusInfo).Methods("GET")

	r.HandleFunc("/tours/{id}/cost", h.UpdateCost).Methods("PUT")

	//tour execution

	r.HandleFunc("/tour-executions/start", executionHandler.StartTour).Methods("POST")
	r.HandleFunc("/tour-executions/{id}/finish", executionHandler.FinishTour).Methods("PUT")
	r.HandleFunc("/tour-executions/{id}/abandon", executionHandler.AbandonTour).Methods("PUT")
	r.HandleFunc("/tour-executions/{id}/update-location", executionHandler.UpdateLocation).Methods("PUT")
	r.HandleFunc("/tour-executions/{id}/add-keypoint", executionHandler.AddCompletedKeyPoint).Methods("PUT")
	r.HandleFunc("/tour-executions/{id}", executionHandler.GetExecutionById).Methods("GET")

	log.Println("Server running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
