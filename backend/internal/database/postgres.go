package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresDB handles PostgreSQL database operations
type PostgresDB struct {
	db *gorm.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB() (*PostgresDB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(
		&User{},
		&SavedRoute{},
		&RoutePreference{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return &PostgresDB{
		db: db,
	}, nil
}

// User represents a user in the system
type User struct {
	gorm.Model
	Email           string           `gorm:"uniqueIndex;not null"`
	Name            string           `gorm:"not null"`
	SavedRoutes     []SavedRoute     `gorm:"foreignKey:UserID"`
	RoutePreference RoutePreference  `gorm:"foreignKey:UserID"`
}

// SavedRoute represents a saved route in the system
type SavedRoute struct {
	gorm.Model
	UserID         uint    `gorm:"not null"`
	StartLat       float64 `gorm:"not null"`
	StartLng       float64 `gorm:"not null"`
	EndLat         float64 `gorm:"not null"`
	EndLng         float64 `gorm:"not null"`
	Distance       float64 `gorm:"not null"` // in meters
	Duration       int64   `gorm:"not null"` // in seconds
	CO2Emission    float64 `gorm:"not null"` // in grams
	TransportMode  string  `gorm:"not null"`
	StartAddress   string
	EndAddress     string
}

// RoutePreference represents user preferences for route calculation
type RoutePreference struct {
	gorm.Model
	UserID             uint    `gorm:"uniqueIndex;not null"`
	PreferredModes     string  `gorm:"not null"` // Comma-separated list
	AvoidHighways      bool    `gorm:"not null"`
	MaxWalkingDistance float64 `gorm:"not null"` // in meters
	PrioritizeEmission bool    `gorm:"not null"`
	MaxTransfers       int     `gorm:"not null"`
}

// CreateUser creates a new user in the database
func (db *PostgresDB) CreateUser(user *User) error {
	return db.db.Create(user).Error
}

// GetUser retrieves a user by ID
func (db *PostgresDB) GetUser(id uint) (*User, error) {
	var user User
	if err := db.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// SaveRoute saves a route to the database
func (db *PostgresDB) SaveRoute(route *SavedRoute) error {
	return db.db.Create(route).Error
}

// GetUserRoutes retrieves all routes for a user
func (db *PostgresDB) GetUserRoutes(userID uint) ([]SavedRoute, error) {
	var routes []SavedRoute
	if err := db.db.Where("user_id = ?", userID).Find(&routes).Error; err != nil {
		return nil, err
	}
	return routes, nil
}

// UpdateRoutePreference updates a user's route preferences
func (db *PostgresDB) UpdateRoutePreference(pref *RoutePreference) error {
	return db.db.Save(pref).Error
}

// GetRoutePreference retrieves a user's route preferences
func (db *PostgresDB) GetRoutePreference(userID uint) (*RoutePreference, error) {
	var pref RoutePreference
	if err := db.db.Where("user_id = ?", userID).First(&pref).Error; err != nil {
		return nil, err
	}
	return &pref, nil
}
