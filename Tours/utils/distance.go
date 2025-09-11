package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
	"tours/model"
)

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

/*
func CalculateTourLength(keypoints []model.KeyPoint) float64 {
	if len(keypoints) < 2 {
		return 0
	}

	var length float64
	for i := 0; i < len(keypoints)-1; i++ {
		length += Haversine(
			keypoints[i].Latitude, keypoints[i].Longitude,
			keypoints[i+1].Latitude, keypoints[i+1].Longitude,
		)
	}
	return length
}
*/

// OSRM response structures
type OSRMResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
	} `json:"routes"`
}

// RoadDistanceCalculator interface for different services
type RoadDistanceCalculator interface {
	CalculateDistance(lat1, lon1, lat2, lon2 float64) (float64, error)
}

// OSRMCalculator implements road distance calculation using OSRM (free public instance)
type OSRMCalculator struct {
	Client *http.Client
}

func NewOSRMCalculator() *OSRMCalculator {
	return &OSRMCalculator{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (osrm *OSRMCalculator) CalculateDistance(lat1, lon1, lat2, lon2 float64) (float64, error) {
	// OSRM public demo server
	baseURL := "http://router.project-osrm.org/route/v1/driving"

	// Format coordinates as lon,lat
	coords := fmt.Sprintf("%.6f,%.6f;%.6f,%.6f", lon1, lat1, lon2, lat2)
	reqURL := fmt.Sprintf("%s/%s?overview=false&alternatives=false&steps=false", baseURL, coords)

	resp, err := osrm.Client.Get(reqURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var osrmResp OSRMResponse
	if err := json.Unmarshal(body, &osrmResp); err != nil {
		return 0, err
	}

	if len(osrmResp.Routes) == 0 {
		return 0, fmt.Errorf("no routes found")
	}

	// Convert meters to kilometers
	return osrmResp.Routes[0].Distance / 1000.0, nil
}

// FallbackCalculator uses haversine if road distance fails
type FallbackCalculator struct {
	Primary RoadDistanceCalculator
}

func NewFallbackCalculator(primary RoadDistanceCalculator) *FallbackCalculator {
	return &FallbackCalculator{Primary: primary}
}

func (f *FallbackCalculator) CalculateDistance(lat1, lon1, lat2, lon2 float64) (float64, error) {
	distance, err := f.Primary.CalculateDistance(lat1, lon1, lat2, lon2)
	if err != nil {
		// Fallback to haversine distance
		return Haversine(lat1, lon1, lat2, lon2), nil
	}
	return distance, nil
}

// Default calculator instance (singleton)
var defaultRoadCalculator RoadDistanceCalculator

// Initialize default calculator
func init() {
	// Use OSRM with fallback to haversine as default
	osrmCalc := NewOSRMCalculator()
	defaultRoadCalculator = NewFallbackCalculator(osrmCalc)
}

// Updated tour length calculation functions
func CalculateTourLength(keypoints []model.KeyPoint) float64 {
	// Use road distance by default now, with fallback to haversine
	return calculateTourLengthWithCalculator(keypoints, defaultRoadCalculator)
}

// For backward compatibility - pure haversine calculation
func CalculateTourLengthHaversine(keypoints []model.KeyPoint) float64 {
	return calculateTourLengthWithCalculator(keypoints, nil)
}

func CalculateTourLengthRoad(keypoints []model.KeyPoint, calculator RoadDistanceCalculator) float64 {
	return calculateTourLengthWithCalculator(keypoints, calculator)
}

func calculateTourLengthWithCalculator(keypoints []model.KeyPoint, calculator RoadDistanceCalculator) float64 {
	if len(keypoints) < 2 {
		return 0
	}

	var length float64
	for i := 0; i < len(keypoints)-1; i++ {
		var distance float64

		if calculator != nil {
			dist, err := calculator.CalculateDistance(
				keypoints[i].Latitude, keypoints[i].Longitude,
				keypoints[i+1].Latitude, keypoints[i+1].Longitude,
			)
			if err != nil {
				// Fallback to haversine on error
				distance = Haversine(
					keypoints[i].Latitude, keypoints[i].Longitude,
					keypoints[i+1].Latitude, keypoints[i+1].Longitude,
				)
			} else {
				distance = dist
			}
		} else {
			distance = Haversine(
				keypoints[i].Latitude, keypoints[i].Longitude,
				keypoints[i+1].Latitude, keypoints[i+1].Longitude,
			)
		}

		length += distance
	}
	return length
}

// Batch calculation for better performance (where supported)
func CalculateTourLengthRoadBatch(keypoints []model.KeyPoint, calculator RoadDistanceCalculator) float64 {
	// For now, fall back to individual calculations
	// Some services support batch requests which would be more efficient
	return CalculateTourLengthRoad(keypoints, calculator)
}
