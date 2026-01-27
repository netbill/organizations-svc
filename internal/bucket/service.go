package bucket

import (
	"context"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	s3 awsxs3
}

func New(awsx3 awsxs3) Bucket {
	return Bucket{
		s3: awsx3,
	}
}

type awsxs3 interface {
	PresignPut(
		ctx context.Context,
		key string,
		ttl time.Duration,
	) (uploadURL, getUrl string, error error)

	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	GetObjectRange(ctx context.Context, key string, maxBytes int64) (io.ReadCloser, error)
	HeadObject(ctx context.Context, key string) (*s3.HeadObjectOutput, error)
	CopyObject(ctx context.Context, tmplKey, finalKey string) (string, error)
	DeleteObject(ctx context.Context, key string) error
}
