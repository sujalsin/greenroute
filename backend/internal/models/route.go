package models

import (
	"time"
)

// TransportMode represents different modes of transportation
type TransportMode string

const (
	Car         TransportMode = "car"
	Bicycle     TransportMode = "bicycle"
	PublicTransit TransportMode = "public_transit"
	Walking     TransportMode = "walking"
)

// Location represents a geographical point
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address,omitempty"`
}

// RouteSegment represents a portion of the route with specific transport mode
type RouteSegment struct {
	StartLocation Location     `json:"start_location"`
	EndLocation   Location     `json:"end_location"`
	Mode          TransportMode `json:"mode"`
	Duration      time.Duration `json:"duration"`
	Distance      float64      `json:"distance"` // in meters
	CO2Emission   float64      `json:"co2_emission"` // in grams
}

// Route represents a complete route with multiple segments
type Route struct {
	ID            string         `json:"id"`
	UserID        string         `json:"user_id"`
	StartLocation Location       `json:"start_location"`
	EndLocation   Location       `json:"end_location"`
	Segments      []RouteSegment `json:"segments"`
	TotalDistance float64        `json:"total_distance"` // in meters
	TotalDuration time.Duration  `json:"total_duration"`
	TotalEmission float64        `json:"total_emission"` // in grams
	CreatedAt     time.Time      `json:"created_at"`
}

// RoutePreferences represents user preferences for route calculation
type RoutePreferences struct {
	PreferredModes    []TransportMode `json:"preferred_modes"`
	AvoidHighways     bool           `json:"avoid_highways"`
	MaxWalkingDistance float64        `json:"max_walking_distance"` // in meters
	PrioritizeEmission bool           `json:"prioritize_emission"`
	MaxTransfers      int            `json:"max_transfers"`
}
