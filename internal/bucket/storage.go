package bucket

import (
	"context"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/netbill/awsx"
)

type bucket interface {
	PresignPut(ctx context.Context, key string, ttl time.Duration) (uploadURL, getURL string, err error)

	HeadObject(ctx context.Context, key string) (*s3.HeadObjectOutput, error)
	GetObject(ctx context.Context, key string) (*s3.GetObjectOutput, error)
	GetObjectRange(ctx context.Context, key string, bytes int64) (*s3.GetObjectOutput, error)

	CopyObject(ctx context.Context, fromKey, toKey string) error
	DeleteObject(ctx context.Context, key string) error
}

type Config struct {
	LinkTTL   time.Duration
	OrgIcon   awsx.ImageValidator
	OrgBanner awsx.ImageValidator
}

type Storage struct {
	s3     bucket
	config Config
}

func NewStorage(s3 bucket, config Config) *Storage {
	return &Storage{
		s3:     s3,
		config: config,
	}
}

func ptrStrEq(a, b *string) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && *a == *b)
}
