package engine

// Engine handles the main functions
type Engine struct {
	redisClient   RedisClient
	sqlClient     SQLClient
	s3Client      S3Client
	googleClient  GoogleClient
	allowedDomain string
}

// Options are meant to be passed to the New function
type Options struct {
	AllowedDomain string
}

// New creates the engine instance
func New(redisClient RedisClient, sqlClient SQLClient, s3Client S3Client, googleClient GoogleClient, options Options) *Engine {
	return &Engine{
		redisClient:   redisClient,
		sqlClient:     sqlClient,
		s3Client:      s3Client,
		googleClient:  googleClient,
		allowedDomain: options.AllowedDomain,
	}
}
