package boot

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/netbill/awsx"
)

type AwsConfig struct {
	BucketName      string `mapstructure:"bucket_name"`
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
}

func newAws(c AwsConfig) *awsx.Bucket {
	s3Client := s3.NewFromConfig(aws.Config{
		Region: c.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			c.AccessKeyID,
			c.SecretAccessKey,
			"",
		),
	})

	return awsx.New(c.BucketName, s3Client, s3.NewPresignClient(s3Client))
}
