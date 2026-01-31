package bucket

import (
	"context"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"
)

const probeBytes = int64(128 * 1024)

type Bucket struct {
	s3                 s3bucket
	OrgIconValidator   ObjectValidator
	OrgBannerValidator ObjectValidator
	tokensTTL          UploadTokensTTL
}

type UploadTokensTTL struct {
	Org time.Duration
}

type Config struct {
	S3                 s3bucket
	OrgIconValidator   ObjectValidator
	OrgBannerValidator ObjectValidator
	UploadTokensTTL    UploadTokensTTL
}

func New(config Config) Bucket {
	return Bucket{
		s3:                 config.S3,
		OrgIconValidator:   config.OrgIconValidator,
		OrgBannerValidator: config.OrgBannerValidator,
	}
}

type s3bucket interface {
	PresignPut(
		ctx context.Context,
		key string,
		ttl time.Duration,
	) (uploadURL, getUrl string, error error)

	GetObjectRange(ctx context.Context, key string, bytes int64) (io.ReadCloser, int64, error)
	CopyObject(ctx context.Context, fromKey, toKey string) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

type ObjectValidator interface {
	ValidateImageResolution(data []byte) (bool, error)
	ValidateImageFormat(data []byte) (bool, error)
	ValidateImageContentType(data []byte) (bool, error)
	ValidateImageSize(size uint) (bool, error)
}
