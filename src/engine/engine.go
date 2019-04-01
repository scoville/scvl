package engine

// Engine handles the main functions
type Engine struct {
	redisClient   RedisClient
	sqlClient     SQLClient
	awsClient     AWSClient
	googleClient  GoogleClient
	allowedDomain string
}

// Options are meant to be passed to the New function
type Options struct {
	AllowedDomain string
}

// New creates the engine instance
func New(redisClient RedisClient, sqlClient SQLClient, awsClient AWSClient, googleClient GoogleClient, options Options) *Engine {
	return &Engine{
		redisClient:   redisClient,
		sqlClient:     sqlClient,
		awsClient:     awsClient,
		googleClient:  googleClient,
		allowedDomain: options.AllowedDomain,
	}
}
