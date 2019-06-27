package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/scoville/scvl/src/adapters/web"
	"github.com/scoville/scvl/src/engine"
	"github.com/scoville/scvl/src/providers/awsclient"
	"github.com/scoville/scvl/src/providers/google"
	"github.com/scoville/scvl/src/providers/redis"
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

	awsClient, err := awsclient.NewClient(awsclient.Config{
		AccessKey:      os.Getenv("AWS_ACCESS_KEY"),
		AccessSecret:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Bucket:       os.Getenv("S3_BUCKET"),
		S3Region:       os.Getenv("S3_REGION"),
		SESRegion:      os.Getenv("SES_REGION"),
		MailFrom:       os.Getenv("MAIL_FROM"),
		MailBCCAddress: os.Getenv("MAIL_BCC_ADDRESS"),
		MainDomain:     os.Getenv("MAIN_DOMAIN"),
		FileDomain:     os.Getenv("FILE_DOMAIN"),
	})
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
		awsClient,
		googleClient,
		engine.Options{
			AllowedDomain: os.Getenv("ALLOWED_DOMAIN"),
		},
	)

	web.Digest = revision
	w := web.New(engine, os.Getenv("SESSION_SECRET"), os.Getenv("MAIN_DOMAIN"))
	log.Fatal(w.Start(":8080"))
}
