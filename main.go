package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/scoville/scvl/src/adapters/web"
	"github.com/scoville/scvl/src/engine"
	"github.com/scoville/scvl/src/providers/google"
	"github.com/scoville/scvl/src/providers/redis"
	"github.com/scoville/scvl/src/providers/s3"
	"github.com/scoville/scvl/src/providers/sql"
)

// build flags
var (
	revision string
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClient := google.NewClient(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_REDIRECT_URL"),
	)

	redisClient, err := redis.NewClient()
	if err != nil {
		log.Fatalf("Failed to create redisClient: %v", err)
	}
	defer redisClient.Close()

	s3Client, err := s3.NewClient(
		os.Getenv("S3_BUCKET"),
		os.Getenv("S3_REGION"),
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
	)
	if err != nil {
		log.Fatalf("Failed to create s3Client: %v", err)
	}

	sqlClient, err := sql.NewClient(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to create sqlClient: %v", err)
	}
	defer sqlClient.Close()

	engine := engine.New(
		redisClient,
		sqlClient,
		s3Client,
		googleClient,
		engine.Options{
			AllowedDomain: os.Getenv("ALLOWED_DOMAIN"),
		},
	)

	web.Digest = revision
	w := web.New(engine, os.Getenv("SESSION_SECRET"))
	log.Fatal(w.Start(":8080"))
}
