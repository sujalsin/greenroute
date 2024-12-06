package external

import (
	"context"
	"errors"
	"fmt"
	"greenroute/internal/models"
	"os"
	"time"

	"googlemaps.github.io/maps"
)

// MapsClient handles interactions with external mapping APIs
type MapsClient struct {
	client *maps.Client
}

// NewMapsClient creates a new instance of MapsClient
func NewMapsClient() (*MapsClient, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return nil, errors.New("Google Maps API key not found in environment variables")
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Maps client: %v", err)
	}

	return &MapsClient{
		client: client,
	}, nil
}

// GetRoute calculates a route between two points using specified transport mode
func (m *MapsClient) GetRoute(
	ctx context.Context,
	origin models.Location,
	destination models.Location,
	mode models.TransportMode,
) (*models.RouteSegment, error) {
	// Convert our transport mode to Google Maps mode
	tMode := convertTransportMode(mode)

	r := &maps.DirectionsRequest{
		Origin:        formatLocation(origin),
		Destination:   formatLocation(destination),
		Mode:          tMode,
		DepartureTime: "now",
		Alternatives:  true,
	}

	routes, _, err := m.client.Directions(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to get directions: %v", err)
	}

	if len(routes) == 0 {
		return nil, errors.New("no routes found")
	}

	// Use the first route
	route := routes[0]
	if len(route.Legs) == 0 {
		return nil, errors.New("route has no legs")
	}

	// Calculate total distance and duration
	leg := route.Legs[0]
	return &models.RouteSegment{
		StartLocation: origin,
		EndLocation:  destination,
		Mode:         mode,
		Duration:     time.Duration(leg.Duration) * time.Second,
		Distance:     float64(leg.Distance.Meters),
		CO2Emission:  calculateEmissions(mode, float64(leg.Distance.Meters)),
	}, nil
}

// formatLocation converts our Location model to Google Maps format
func formatLocation(loc models.Location) string {
	return fmt.Sprintf("%f,%f", loc.Latitude, loc.Longitude)
}

// convertTransportMode converts our transport mode to Google Maps format
func convertTransportMode(mode models.TransportMode) maps.Mode {
	switch mode {
	case models.Car:
		return maps.TravelModeDriving
	case models.Bicycle:
		return maps.TravelModeBicycling
	case models.Walking:
		return maps.TravelModeWalking
	case models.PublicTransit:
		return maps.TravelModeTransit
	default:
		return maps.TravelModeDriving
	}
}

// calculateEmissions estimates CO2 emissions based on transport mode and distance
func calculateEmissions(mode models.TransportMode, distanceMeters float64) float64 {
	// Emissions in grams of CO2 per kilometer
	emissionsPerKm := map[models.TransportMode]float64{
		models.Car:         120.0, // Average car
		models.Bicycle:     0.0,   // Zero emissions
		models.Walking:     0.0,   // Zero emissions
		models.PublicTransit: 60.0, // Average bus/train
	}

	kmTraveled := distanceMeters / 1000.0
	return kmTraveled * emissionsPerKm[mode]
}
