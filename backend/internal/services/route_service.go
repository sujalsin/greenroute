package services

import (
	"context"
	"errors"
	"greenroute/internal/database"
	"greenroute/internal/external"
	"greenroute/internal/models"
	"time"
)

// RouteService handles route calculation and optimization
type RouteService struct {
	mapsClient     *external.MapsClient
	chargingClient *external.ChargingClient
	postgres       *database.PostgresDB
	mongodb        *database.MongoDB
}

// NewRouteService creates a new instance of RouteService
func NewRouteService(
	mapsClient *external.MapsClient,
	chargingClient *external.ChargingClient,
	postgres *database.PostgresDB,
	mongodb *database.MongoDB,
) *RouteService {
	return &RouteService{
		mapsClient:     mapsClient,
		chargingClient: chargingClient,
		postgres:       postgres,
		mongodb:        mongodb,
	}
}

// RouteWithCharging represents a route with EV charging stations
type RouteWithCharging struct {
	Route           *models.Route
	ChargingStations []external.ChargingStation
}

// CalculateRoute generates an optimized route based on user preferences
func (s *RouteService) CalculateRoute(
	ctx context.Context,
	start models.Location,
	end models.Location,
	prefs models.RoutePreferences,
	userID string,
) (*RouteWithCharging, error) {
	if !s.validateLocations(start, end) {
		return nil, errors.New("invalid locations provided")
	}

	// Get historical traffic data
	now := time.Now()
	pattern, err := s.mongodb.GetTrafficPattern(
		start.Latitude, start.Longitude,
		end.Latitude, end.Longitude,
		int(now.Weekday()),
		now.Hour(),
	)
	if err != nil {
		return nil, err
	}

	// Calculate routes for each preferred mode
	var segments []models.RouteSegment
	var totalDistance, totalEmission float64
	var totalDuration time.Duration
	var waypoints []struct{ Lat, Lng float64 }

	// Add start point to waypoints
	waypoints = append(waypoints, struct{ Lat, Lng float64 }{
		Lat: start.Latitude,
		Lng: start.Longitude,
	})

	for _, mode := range prefs.PreferredModes {
		segment, err := s.mapsClient.GetRoute(ctx, start, end, mode)
		if err != nil {
			continue // Skip this mode if calculation fails
		}

		// Adjust duration based on historical data if available
		if pattern != nil {
			segment.Duration = time.Duration(pattern.Duration) * time.Second
		}

		segments = append(segments, *segment)
		totalDistance += segment.Distance
		totalEmission += segment.CO2Emission
		totalDuration += segment.Duration
	}

	if len(segments) == 0 {
		return nil, errors.New("no valid routes found for any preferred mode")
	}

	// Add end point to waypoints
	waypoints = append(waypoints, struct{ Lat, Lng float64 }{
		Lat: end.Latitude,
		Lng: end.Longitude,
	})

	// Create the complete route
	route := &models.Route{
		UserID:        userID,
		StartLocation: start,
		EndLocation:   end,
		Segments:      segments,
		TotalDistance: totalDistance,
		TotalDuration: totalDuration,
		TotalEmission: totalEmission,
		CreatedAt:     time.Now(),
	}

	// Save the route for future reference
	if err := s.saveRoute(route); err != nil {
		return nil, err
	}

	// Update traffic pattern
	s.updateTrafficPattern(start, end, segments[0].Duration)

	// Find charging stations along the route
	stations, err := s.chargingClient.FindStationsAlongRoute(waypoints, 2.0) // 2km corridor
	if err != nil {
		// Don't fail the request if charging station lookup fails
		stations = []external.ChargingStation{}
	}

	return &RouteWithCharging{
		Route:           route,
		ChargingStations: stations,
	}, nil
}

// validateLocations checks if the provided locations are valid
func (s *RouteService) validateLocations(start, end models.Location) bool {
	return isValidLatitude(start.Latitude) &&
		isValidLatitude(end.Latitude) &&
		isValidLongitude(start.Longitude) &&
		isValidLongitude(end.Longitude)
}

func isValidLatitude(lat float64) bool {
	return lat >= -90 && lat <= 90
}

func isValidLongitude(lon float64) bool {
	return lon >= -180 && lon <= 180
}

// saveRoute saves the route to PostgreSQL
func (s *RouteService) saveRoute(route *models.Route) error {
	savedRoute := &database.SavedRoute{
		UserID:        route.UserID,
		StartLat:      route.StartLocation.Latitude,
		StartLng:      route.StartLocation.Longitude,
		EndLat:        route.EndLocation.Latitude,
		EndLng:        route.EndLocation.Longitude,
		Distance:      route.TotalDistance,
		Duration:      int64(route.TotalDuration.Seconds()),
		CO2Emission:   route.TotalEmission,
		TransportMode: string(route.Segments[0].Mode), // Use the primary mode
		StartAddress:  route.StartLocation.Address,
		EndAddress:    route.EndLocation.Address,
	}

	return s.postgres.SaveRoute(savedRoute)
}

// updateTrafficPattern updates the traffic pattern in MongoDB
func (s *RouteService) updateTrafficPattern(start, end models.Location, duration time.Duration) {
	now := time.Now()
	pattern := &database.TrafficPattern{
		StartLat:    start.Latitude,
		StartLng:    start.Longitude,
		EndLat:      end.Latitude,
		EndLng:      end.Longitude,
		DayOfWeek:   int(now.Weekday()),
		HourOfDay:   now.Hour(),
		Duration:    duration.Seconds(),
		Timestamp:   now,
		SampleCount: 1,
	}

	// Ignore error as this is not critical
	_ = s.mongodb.SaveTrafficPattern(pattern)
}
