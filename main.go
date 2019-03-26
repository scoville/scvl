package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/scoville/scvl/src/adapters/web"
)

// build flags
var (
	revision string
)

var client *redisClient
var s3Client *S3Client

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	web.Digest = revision

	sessionSecret := os.Getenv("SESSION_SECRET")

	client, err = newRedisClient()
	if err != nil {
		log.Fatalf("Failed to create redisClient: %v", err)
	}
	defer client.Close()

	s3Client, err = newS3Client(os.Getenv("S3_BUCKET"), os.Getenv("S3_REGION"), os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_ACCESS_KEY"))
	if err != nil {
		log.Fatalf("Failed to create s3Client: %v", err)
	}

	setupGoogleConfig()

	dbURL := os.Getenv("DB_URL")

	setupManager()
	defer manager.db.Close()
	log.Fatal()
}
