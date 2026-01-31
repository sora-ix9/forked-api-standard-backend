package s3

import (
	"context"
	"fdlp-standard-api/pkg/config"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewClient initializes and returns a new AWS S3 client
func NewClient(cfg *config.Config) *s3.Client {
	var awsCfg aws.Config
	var err error

	if cfg.AWSAccessKeyID != "" && cfg.AWSSecretAccessKey != "" {
		creds := credentials.NewStaticCredentialsProvider(cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, "")
		awsCfg, err = awsconfig.LoadDefaultConfig(context.Background(),
			awsconfig.WithRegion(cfg.AWSRegion),
			awsconfig.WithCredentialsProvider(creds),
		)
	} else {
		// Fallback to default credential chain (e.g. environment variables, IAM role)
		awsCfg, err = awsconfig.LoadDefaultConfig(context.Background(),
			awsconfig.WithRegion(cfg.AWSRegion),
		)
	}

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return s3.NewFromConfig(awsCfg)
}
