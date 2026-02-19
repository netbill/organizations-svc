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
		BucketName      string
		Region          string
		AccessKeyID     string
		SecretAccessKey string
		SessionToken    string
	}
	Media struct {
		Link struct {
			TTL time.Duration
		}
		Organization struct {
			Icon   awsx.ImageValidator
			Banner awsx.ImageValidator
		}
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
