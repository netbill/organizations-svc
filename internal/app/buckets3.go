package app

import (
	"context"
	"fmt"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/netbill/awsx"
	"github.com/netbill/organizations-svc/internal/bucket"
)

func (a *App) BuildBucket() bucket.Storage {
	cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
		awscfg.WithRegion(a.config.S3.Aws.Region),
		awscfg.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				a.config.S3.Aws.AccessKeyID,
				a.config.S3.Aws.SecretAccessKey,
				a.config.S3.Aws.SessionToken,
			),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to load S3 config: %v", err))
	}

	return bucket.NewStorage(awsx.New(a.config.S3.Aws.BucketName, cfg), bucket.Config{
		LinkTTL:   a.config.S3.Media.Link.TTL,
		OrgIcon:   a.config.S3.Media.Organization.Icon,
		OrgBanner: a.config.S3.Media.Organization.Banner,
	})
}
