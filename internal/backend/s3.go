package backend

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Config holds S3 backend configuration.
type S3Config struct {
	Bucket string
	Prefix string
	Region string
}

// S3 implements Backend for S3-compatible storage.
type S3 struct {
	client *s3.Client
	bucket string
	prefix string
}

// NewS3 creates an S3 backend. Uses default credential chain (env, shared config).
func NewS3(ctx context.Context, cfg S3Config) (*S3, error) {
	opts := []func(*config.LoadOptions) error{}
	if cfg.Region != "" {
		opts = append(opts, config.WithRegion(cfg.Region))
	}
	awsCfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}
	client := s3.NewFromConfig(awsCfg)
	return &S3{client: client, bucket: cfg.Bucket, prefix: cfg.Prefix}, nil
}

func (s *S3) key(k string) string {
	if s.prefix == "" {
		return k
	}
	return s.prefix + "/" + k
}

// Put implements Backend.
func (s *S3) Put(key string, r io.Reader) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key(key)),
		Body:   bytes.NewReader(body),
	})
	return err
}

// Get implements Backend.
func (s *S3) Get(key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.key(key)),
	})
	if err != nil {
		return nil, err
	}
	return out.Body, nil
}

// List implements Backend. Returns keys without the backend prefix.
func (s *S3) List(prefix string) ([]string, error) {
	fullPrefix := s.key(prefix)
	strip := len(s.prefix) + 1
	if s.prefix == "" {
		strip = 0
	}
	var keys []string
	pager := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(fullPrefix),
	})
	for pager.HasMorePages() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			if obj.Key != nil {
				k := *obj.Key
				if strip > 0 && len(k) >= strip {
					k = k[strip:]
				}
				keys = append(keys, k)
			}
		}
	}
	return keys, nil
}
