// Package imagestore is the S3 client for the Image Store (ADR-0004):
// Cloudflare R2 through aws-sdk-go-v2, configured as
// docs/research/cloudflare-r2-and-deploy.md §1 says (region auto, custom
// endpoint, checksums only when required). The builder uploads with it;
// the prune command lists and deletes with the same client and the same
// environment.
package imagestore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"golang.org/x/sync/errgroup"
)

// Environment variables holding the credentials. They never live in
// site.yaml.
const (
	EnvAccountID       = "CLOUDFLARE_ACCOUNT_ID"
	EnvAccessKeyID     = "R2_ACCESS_KEY_ID"
	EnvSecretAccessKey = "R2_SECRET_ACCESS_KEY"
	// EnvEndpoint optionally overrides the S3 endpoint (default
	// https://<account id>.r2.cloudflarestorage.com): a jurisdiction
	// endpoint (.eu.) or a test fake.
	EnvEndpoint = "R2_ENDPOINT"
)

// CacheControl is set on every object: keys are content-hashed, so an
// object never changes and caches for a year. The builder sends the same
// header for /assets/ on both Targets.
const CacheControl = "public, max-age=31536000, immutable"

// Credentials are the R2 account and API token (Access Key ID = the
// token's id, Secret = SHA-256 of its value), plus the endpoint.
type Credentials struct {
	AccountID, AccessKeyID, SecretAccessKey string
	// Endpoint is the S3 endpoint URL; FromEnv defaults it to
	// https://<AccountID>.r2.cloudflarestorage.com.
	Endpoint string
}

// FromEnv reads Credentials through getenv. A missing variable is an error
// naming every one that is unset.
func FromEnv(getenv func(string) string) (Credentials, error) {
	c := Credentials{
		AccountID:       getenv(EnvAccountID),
		AccessKeyID:     getenv(EnvAccessKeyID),
		SecretAccessKey: getenv(EnvSecretAccessKey),
		Endpoint:        getenv(EnvEndpoint),
	}
	var missing []string
	for _, v := range []struct{ name, value string }{
		{EnvAccountID, c.AccountID}, {EnvAccessKeyID, c.AccessKeyID}, {EnvSecretAccessKey, c.SecretAccessKey},
	} {
		if v.value == "" {
			missing = append(missing, v.name)
		}
	}
	if len(missing) > 0 {
		noun := "variable"
		if len(missing) > 1 {
			noun = "variables"
		}
		return Credentials{}, fmt.Errorf("environment %s %s not set", noun, strings.Join(missing, ", "))
	}
	if c.Endpoint == "" {
		c.Endpoint = "https://" + c.AccountID + ".r2.cloudflarestorage.com"
	}
	return c, nil
}

// Client talks to one bucket.
type Client struct {
	s3     *s3.Client
	bucket string
}

// New returns a client for bucket. It is built from explicit options, not
// config.LoadDefaultConfig, so nothing in ~/.aws or AWS_* leaks in.
func New(c Credentials, bucket string) *Client {
	return &Client{bucket: bucket, s3: s3.New(s3.Options{
		Region:       "auto", // R2's only region; required by the SDK, unused by R2
		Credentials:  credentials.NewStaticCredentialsProvider(c.AccessKeyID, c.SecretAccessKey, ""),
		BaseEndpoint: aws.String(c.Endpoint),
		// Bucket in the path, not the host, so any endpoint (a jurisdiction
		// one, a test fake on 127.0.0.1) works unchanged.
		UsePathStyle: true,
		// The SDK's default CRC32 on every PUT is not something R2 documents
		// accepting on single PUTs (research §1.3).
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
	})}
}

// Object is one image to upload. Data, if nil, is read from File.
type Object struct {
	Key         string
	ContentType string
	Data        []byte
	File        string
}

// Bytes is the object's content: Data, or File's bytes when Data is nil.
func (o Object) Bytes() ([]byte, error) {
	if o.Data != nil {
		return o.Data, nil
	}
	return os.ReadFile(o.File)
}

// Upload makes every object exist in the bucket: HEAD each key and PUT only
// on a 404, with ContentType and CacheControl, on at most workers
// goroutines. Keys are content-hashed, so an existing key already holds
// the same bytes. The first failure (non-2xx, unreadable file) stops the
// upload and is returned naming the key.
func (c *Client) Upload(ctx context.Context, objs []Object, workers int) error {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(max(1, workers))
	for _, o := range objs {
		g.Go(func() error {
			if err := c.ensure(ctx, o); err != nil {
				return fmt.Errorf("upload %s: %w", o.Key, err)
			}
			return nil
		})
	}
	return g.Wait()
}

func (c *Client) ensure(ctx context.Context, o Object) error {
	_, err := c.s3.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &c.bucket, Key: &o.Key})
	if err == nil {
		return nil
	}
	if nf := (*types.NotFound)(nil); !errors.As(err, &nf) {
		return err
	}
	data, err := o.Bytes()
	if err != nil {
		return err
	}
	_, err = c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       &c.bucket,
		Key:          &o.Key,
		Body:         bytes.NewReader(data),
		ContentType:  aws.String(o.ContentType),
		CacheControl: aws.String(CacheControl),
	})
	return err
}
