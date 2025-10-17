package config

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JwtSecret []byte
var jwtKey = []byte("your-secret-key")

func InitDB() {
	// Get database connection details from environment variables with fallbacks
	// dbHost := getEnv("DB_HOST", "localhost")
	// dbUser := getEnv("DB_USER", "postgres")
	// dbPassword := getEnv("DB_PASSWORD", "123")
	// dbName := getEnv("DB_NAME", "gym_test")
	// dbPort := getEnv("DB_PORT", "5432")


	// Session Pooler settings (IPv4 compatible)
	dbHost := getEnv("DB_HOST", "aws-1-us-east-2.pooler.supabase.com")
	dbUser := getEnv("DB_USER", "postgres.totqgnazgprwdvwdqmpe") // Supabase Session Pooler requires full username
	dbPassword := getEnv("DB_PASSWORD", "1234")
	dbName := getEnv("DB_NAME", "postgres")
	dbPort := getEnv("DB_PORT", "5432")




	// Construct DSN from environment variables
	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + 
	       " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=UTC"

	log.Printf("Connecting to database at %s:%s", dbHost, dbPort)
	
	// Retry connection with backoff
	var err error
	for i := 0; i < 10; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Failed to connect to database (attempt %d/10): %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database after retries:", err)
	}
	
	log.Println("Successfully connected to database")
}

func InitJWT() {
	JwtSecret = []byte(getEnv("JWT_SECRET", "fallback-secret-key-change-in-production"))
}

func GenerateJWT(userID string, userType string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID, // UUID stored as string
		"user_type": userType,
		"exp":       time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}