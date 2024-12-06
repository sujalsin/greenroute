package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// ChargingClient handles interactions with EV charging station APIs
type ChargingClient struct {
	apiKey string
	client *http.Client
}

// ChargingStation represents an EV charging station
type ChargingStation struct {
	ID          int     `json:"ID"`
	AddressInfo struct {
		Title     string  `json:"Title"`
		Address   string  `json:"AddressLine1"`
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
	} `json:"AddressInfo"`
	Connections []struct {
		ConnectionType struct {
			Title string `json:"Title"`
		} `json:"ConnectionType"`
		PowerKW float64 `json:"PowerKW"`
	} `json:"Connections"`
	UsageType struct {
		Title string `json:"Title"`
	} `json:"UsageType"`
}

// NewChargingClient creates a new instance of ChargingClient
func NewChargingClient() (*ChargingClient, error) {
	apiKey := os.Getenv("OPENCHARGE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OpenChargeMap API key not found in environment variables")
	}

	return &ChargingClient{
		apiKey: apiKey,
		client: &http.Client{},
	}, nil
}

// FindNearbyStations finds charging stations near a location within a radius
func (c *ChargingClient) FindNearbyStations(lat, lng float64, radiusKm float64) ([]ChargingStation, error) {
	url := fmt.Sprintf(
		"https://api.openchargemap.io/v3/poi?output=json&latitude=%f&longitude=%f&distance=%f&distanceunit=km&maxresults=10",
		lat, lng, radiusKm,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("X-API-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var stations []ChargingStation
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return stations, nil
}

// FindStationsAlongRoute finds charging stations along a route within a corridor
func (c *ChargingClient) FindStationsAlongRoute(waypoints []struct{ Lat, Lng float64 }, corridorKm float64) ([]ChargingStation, error) {
	var allStations []ChargingStation
	seenStations := make(map[int]bool)

	// Search for stations near each waypoint
	for _, wp := range waypoints {
		stations, err := c.FindNearbyStations(wp.Lat, wp.Lng, corridorKm)
		if err != nil {
			continue
		}

		// Deduplicate stations
		for _, station := range stations {
			if !seenStations[station.ID] {
				allStations = append(allStations, station)
				seenStations[station.ID] = true
			}
		}
	}

	return allStations, nil
}
