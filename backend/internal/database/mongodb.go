package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB handles MongoDB database operations
type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB() (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db := client.Database("greenroute")
	return &MongoDB{
		client: client,
		db:     db,
	}, nil
}

// TrafficPattern represents historical traffic data for a route segment
type TrafficPattern struct {
	StartLat    float64   `bson:"start_lat"`
	StartLng    float64   `bson:"start_lng"`
	EndLat      float64   `bson:"end_lng"`
	EndLng      float64   `bson:"end_lat"`
	DayOfWeek   int       `bson:"day_of_week"` // 0 = Sunday, 6 = Saturday
	HourOfDay   int       `bson:"hour_of_day"` // 0-23
	Duration    float64   `bson:"duration"`     // Average duration in seconds
	Timestamp   time.Time `bson:"timestamp"`
	SampleCount int       `bson:"sample_count"`
}

// SaveTrafficPattern saves a traffic pattern to MongoDB
func (m *MongoDB) SaveTrafficPattern(pattern *TrafficPattern) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.Collection("traffic_patterns")
	
	// Try to update existing pattern or insert new one
	filter := bson.M{
		"start_lat":   pattern.StartLat,
		"start_lng":   pattern.StartLng,
		"end_lat":     pattern.EndLat,
		"end_lng":     pattern.EndLng,
		"day_of_week": pattern.DayOfWeek,
		"hour_of_day": pattern.HourOfDay,
	}

	update := bson.M{
		"$inc": bson.M{
			"sample_count": 1,
			"duration":     pattern.Duration,
		},
		"$set": bson.M{
			"timestamp": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// GetTrafficPattern retrieves historical traffic data for a route segment
func (m *MongoDB) GetTrafficPattern(
	startLat, startLng, endLat, endLng float64,
	dayOfWeek, hourOfDay int,
) (*TrafficPattern, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := m.db.Collection("traffic_patterns")

	filter := bson.M{
		"start_lat":   startLat,
		"start_lng":   startLng,
		"end_lat":     endLat,
		"end_lng":     endLng,
		"day_of_week": dayOfWeek,
		"hour_of_day": hourOfDay,
	}

	var pattern TrafficPattern
	err := collection.FindOne(ctx, filter).Decode(&pattern)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &pattern, nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}
