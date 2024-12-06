package database

import (
	"time"
)

// SavedRoute represents a route stored in PostgreSQL
type SavedRoute struct {
	ID            string    `gorm:"primarykey"`
	UserID        string    `gorm:"not null"`
	StartLat      float64   `gorm:"not null"`
	StartLng      float64   `gorm:"not null"`
	EndLat        float64   `gorm:"not null"`
	EndLng        float64   `gorm:"not null"`
	Distance      float64   `gorm:"not null"`
	Duration      int64     `gorm:"not null"` // stored in seconds
	CO2Emission   float64   `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
}

// TrafficPattern represents traffic data stored in MongoDB
type TrafficPattern struct {
	ID          string    `bson:"_id,omitempty"`
	StartLat    float64   `bson:"start_lat"`
	StartLng    float64   `bson:"start_lng"`
	EndLat      float64   `bson:"end_lat"`
	EndLng      float64   `bson:"end_lng"`
	Duration    int64     `bson:"duration"` // in seconds
	TimeOfDay   int       `bson:"time_of_day"` // hour of day (0-23)
	DayOfWeek   int       `bson:"day_of_week"` // 0 = Sunday, 6 = Saturday
	UpdatedAt   time.Time `bson:"updated_at"`
}
