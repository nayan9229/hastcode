package imagestore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ObjectInfo is one object a listing returned.
type ObjectInfo struct {
	Key  string
	Size int64
}

// List returns every object in the bucket whose key starts with prefix
// ("" for all), in key order, following ListObjectsV2 continuation
// tokens (up to 1,000 keys per page, one Class A operation each).
func (c *Client) List(ctx context.Context, prefix string) ([]ObjectInfo, error) {
	var objs []ObjectInfo
	p := s3.NewListObjectsV2Paginator(c.s3, &s3.ListObjectsV2Input{Bucket: &c.bucket, Prefix: aws.String(prefix)})
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, o := range page.Contents {
			objs = append(objs, ObjectInfo{Key: aws.ToString(o.Key), Size: aws.ToInt64(o.Size)})
		}
	}
	return objs, nil
}
