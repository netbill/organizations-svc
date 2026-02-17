package bucket

import (
	"context"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/netbill/awsx"
)

type Config struct {
	Aws struct {
		BucketName      string `mapstructure:"bucket_name"`
		Region          string `mapstructure:"region"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		SessionToken    string `mapstructure:"session_token"`
	} `mapstructure:"aws"`

	Media struct {
		Link struct {
			TTL time.Duration `mapstructure:"ttl"`
		} `mapstructure:"link"`
		Organization struct {
			Icon   awsx.ImageValidator `mapstructure:"icon"`
			Banner awsx.ImageValidator `mapstructure:"banner"`
		} `mapstructure:"organization"`
	}
}

type Bucket struct {
	s3     awsx.Bucket
	config Config
}

func New(config Config) (Bucket, error) {
	cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
		awscfg.WithRegion(config.Aws.Region),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.Aws.AccessKeyID,
				config.Aws.SecretAccessKey,
				config.Aws.SessionToken,
			),
		),
	)
	if err != nil {
		return Bucket{}, err
	}

	bucket := awsx.New(config.Aws.BucketName, cfg)

	return Bucket{
		s3:     bucket,
		config: config,
	}, nil
}

func ptrStrEq(a, b *string) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && *a == *b)
}
