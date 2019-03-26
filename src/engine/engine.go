package engine

// Engine handles the main functions
type Engine struct {
	redisClient RedisClient
	sqlClient   SQLClient
	s3Client    S3Client
}

// New creates the engine instance
func New(redisClient RedisClient, sqlClient SQLClient, s3Client S3Client) *Engine {
	return &Engine{
		redisClient: redisClient,
		sqlClient:   sqlClient,
		s3Client:    s3Client,
	}
}
