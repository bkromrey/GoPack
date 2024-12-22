package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// GetGeoAPIKey returns the Geocoding API key string stored in the .env file
func GetGeoAPIKey() string {

	// load environment variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	geoAPI := os.Getenv("GEO_API_KEY")

	return geoAPI
}
