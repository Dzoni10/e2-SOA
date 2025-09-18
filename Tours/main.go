package main

import (
	"log"
	"net/http"
	"tours/database"
	"tours/handler"
	"tours/repo"
	"tours/service"

	"github.com/gorilla/mux"
)

func main() {

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

	/*corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "UPDATE", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})*/

	log.Println("Server running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
