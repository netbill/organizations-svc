package bucket

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

const probeBytes = 128 * 1024
const OrganizationUploadTTL time.Duration = 1 * time.Hour

type Bucket struct {
	s3 s3storage
}

func New(s3 s3storage) Bucket {
	return Bucket{
		s3: s3,
	}
}

type s3storage interface {
	PresignPut(
		ctx context.Context,
		key string,
		ttl time.Duration,
	) (uploadURL, getUrl string, error error)

	HeadObject(ctx context.Context, key string) (*s3.HeadObjectOutput, error)
	GetObjectRange(ctx context.Context, key string, bytes int64) (io.ReadCloser, int64, error)
	CopyObject(ctx context.Context, fromKey, toKey string) (string, error)
	DeleteObject(ctx context.Context, key string) error
}

func allowed(value string, allowed []string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || len(allowed) == 0 {
		return false
	}
	for _, a := range allowed {
		a = strings.ToLower(strings.TrimSpace(a))
		if value == a {
			return true
		}
	}
	return false
}

func (b Bucket) ObjectExists(
	ctx context.Context,
	key string,
) (bool, error) {
	_, err := b.s3.HeadObject(ctx, key)
	if err != nil {
		var respErr *smithyhttp.ResponseError
		if errors.As(err, &respErr) && (respErr.HTTPStatusCode() == 404 || respErr.HTTPStatusCode() == 403) {
			return false, nil
		}

		return false, fmt.Errorf("failed to head object, cause: %w", err)
	}

	return true, nil
}

func (b Bucket) ValidateImgSize(
	ctx context.Context,
	key string,
	maxLength int64,
) (bool, error) {
	obj, err := b.s3.HeadObject(ctx, key)
	if err != nil {
		var respErr *smithyhttp.ResponseError
		if errors.As(err, &respErr) && (respErr.HTTPStatusCode() == 404 || respErr.HTTPStatusCode() == 403) {
			return false, nil
		}

		return false, fmt.Errorf("failed to head object, cause: %w", err)
	}

	if obj.ContentRange == nil {
		return false, fmt.Errorf("content-range header is missing")
	}
	parts := strings.Split(strings.TrimSpace(*obj.ContentRange), "/")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid content-range: %q", *obj.ContentRange)
	}

	size, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid content-range size: %q", parts[1])
	}

	if size > maxLength {
		return false, nil
	}

	return true, nil
}

func (b Bucket) validateImgObjet(
	ctx context.Context,
	key string,
	maxLength int64,
	allowedContentTypes []string,
	allowedContentFormats []string,
	maxW, maxH int64,
) error {
	obj, err := b.s3.HeadObject(ctx, key)
	if err != nil {
		var respErr *smithyhttp.ResponseError
		if errors.As(err, &respErr) && (respErr.HTTPStatusCode() == 404 || respErr.HTTPStatusCode() == 403) {
			return errx.ErrorNoContentUploaded.Raise(
				fmt.Errorf("image upload not found for key %s, cause: %s", key, err),
			)
		}
		return fmt.Errorf("failed to head object, cause: %w", err)
	}

	if obj.ContentRange == nil {
		return fmt.Errorf("content-range header is missing")
	}
	parts := strings.Split(strings.TrimSpace(*obj.ContentRange), "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid content-range: %q", *obj.ContentRange)
	}

	size, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid content-range size: %q", parts[1])
	}

	if size > maxLength {
		return errx.ErrorContentSizeExceed.Raise(
			fmt.Errorf("object content size exceed allowed"),
		)
	}

	reader, err := b.s3.GetObjectRange(ctx, key, probeBytes)
	if err != nil {
		return fmt.Errorf("failed to getting object range, cause: %w", err)
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to reading object probe, cause: %w", err)
	}
	detectedCT := http.DetectContentType(content)
	if len(allowedContentTypes) > 0 && !allowed(detectedCT, allowedContentTypes) {
		return errx.ErrorContentTypeIsNotAllowed.Raise(
			fmt.Errorf("content type is not allowed"),
		)
	}

	img, format, err := image.DecodeConfig(bytes.NewReader(content))
	if err != nil {
		return fmt.Errorf("failed to decode image config: %w", err)
	}
	if len(allowedContentFormats) > 0 && !allowed(format, allowedContentFormats) {
		return errx.ErrorImageFormatIsNotAllowed.Raise(
			fmt.Errorf("organization icon image format is not allowed"),
		)
	}
	if int64(img.Width) > maxW || int64(img.Height) > maxH {
		return errx.ErrorImageResolutionExceed.Raise(
			fmt.Errorf("organization icon image resolution exceed allowed"),
		)
	}

	return nil
}
