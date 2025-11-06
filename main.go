package main

import (
	"fmt"
	"os"
	"strings"

	"api/internal/mongodb"
	"api/internal/plog"
	"api/middleware"
	"api/routers"
	"api/services/minio"
	"api/services/oauth2/google"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// cleanEndpoint removes https:// prefix if it exists
func cleanEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	if strings.HasPrefix(endpoint, "https://") {
		return strings.TrimPrefix(endpoint, "https://")
	}
	return endpoint
}

func main() {
	logger := plog.NewBizLogger("main")

	env := os.Getenv("ENV")

	err := godotenv.Load(os.ExpandEnv("./config/.env"))
	if err != nil {
		err := godotenv.Load(os.ExpandEnv(".env"))
		if err != nil {
			fmt.Println(err)
		} else {
			logger.Info().Msg("loaded environment variables from .env")
		}
	} else {
		logger.Info().Msg("loaded environment variables from ./config/.env")
	}

	_, _, _, err = mongodb.ConnectMongoWithString(
		os.Getenv("MONGODB_URI"),
		os.Getenv("MONGODB_DATABASE"), 100, nil, nil)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.DebugMode)

	logger.Info().Msgf("starting server in %s mode", env)

	r := gin.New()

	// setup gin middleware
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())
	routers.InitRouter(r)

	// setup google oauth2 client
	_, err = google.Config(&google.ConfigOptions{
		ClientID:          os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret:      os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:       os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientRedirectURL: os.Getenv("GOOGLE_AUTH_REDIRECT_URL"),
	})
	if err != nil {
		logger.Error().Msgf("error setting up Google OAuth2 client: %v", err)
	} else {
		logger.Info().Msg("Google OAuth2 client setup successfully")
	}

	// Setup minio client
	minioEndpoint := cleanEndpoint(os.Getenv("MIN_ENDPOINT"))
	logger.Info().Msgf("connecting to MinIO endpoint: %s", minioEndpoint)
	_, err = minio.NewClient(
		minioEndpoint,
		os.Getenv("MIN_ACCESSKEY"),
		os.Getenv("MIN_SECRETKEY"),
		true,
	)
	if err != nil {
		logger.Error().Msgf("error setting up MinIO client: %v", err)
	} else {
		logger.Info().Msg("MinIO client setup successfully")

		// Ensure required buckets exist
		requiredBuckets := []string{"documents", "images", "uploads"}
		for _, bucket := range requiredBuckets {
			err = minio.EnsureBucketExists(bucket)
			if err != nil {
				logger.Error().Msgf("error ensuring bucket %s exists: %v", bucket, err)
			} else {
				logger.Info().Msgf("bucket %s verified/created successfully", bucket)
			}
		}
	}

	logger.Info().Msg("Starting server on " + fmt.Sprintf("%v:%v", os.Getenv("API_HOST"), os.Getenv("API_PORT")))
	if err := r.Run(fmt.Sprintf("%v:%v", os.Getenv("API_HOST"), os.Getenv("API_PORT"))); err != nil {
		logger.Error().Msgf("error starting server: %v", err)
		os.Exit(1)
	}
}
