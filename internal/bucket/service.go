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

type Config struct {
	Link struct {
		TTL time.Duration `json:"ttl"`
	} `json:"link"`
	Organization struct {
		Icon struct {
			AllowedFormats   []string `mapstructure:"allowed_formats" required:"true"`
			MaxWidth         int      `mapstructure:"max_width" required:"true"`
			MaxHeight        int      `mapstructure:"max_height" required:"true"`
			ContentLengthMax int      `mapstructure:"content_length_max" required:"true"`
		} `mapstructure:"icon"`
		Banner struct {
			AllowedFormats   []string `mapstructure:"allowed_formats" required:"true"`
			MaxWidth         int      `mapstructure:"max_width" required:"true"`
			MaxHeight        int      `mapstructure:"max_height" required:"true"`
			ContentLengthMax int      `mapstructure:"content_length_max" required:"true"`
		}
	} `mapstructure:"organization"`
}

type Bucket struct {
	s3     storage
	config Config
}

func New(s3 storage, config Config) Bucket {
	return Bucket{
		s3:     s3,
		config: config,
	}
}

type storage interface {
	PresignPut(
		ctx context.Context,
		key string,
		ttl time.Duration,
	) (uploadURL, getUrl string, error error)

	HeadObject(
		ctx context.Context,
		key string,
	) (*s3.HeadObjectOutput, error)

	GetObjectRange(
		ctx context.Context,
		key string,
		bytes int64,
	) (body io.ReadCloser, err error)

	CopyObject(ctx context.Context, tmplKey, finalKey string) (string, error)
	DeleteObject(ctx context.Context, key string) error
}
