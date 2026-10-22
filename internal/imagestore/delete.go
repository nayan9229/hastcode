package imagestore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// DeleteBatch is the most keys one DeleteObjects request may name (S3 and
// R2: 1,000).
const DeleteBatch = 1000

// Delete removes keys with DeleteObjects, DeleteBatch keys per request.
// It stops at the first failed request, or the first batch in which R2
// reports a key it could not delete, naming that key.
func (c *Client) Delete(ctx context.Context, keys []string) error {
	for start := 0; start < len(keys); start += DeleteBatch {
		batch := keys[start:min(start+DeleteBatch, len(keys))]
		ids := make([]types.ObjectIdentifier, len(batch))
		for i, k := range batch {
			ids[i] = types.ObjectIdentifier{Key: aws.String(k)}
		}
		out, err := c.s3.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: &c.bucket,
			Delete: &types.Delete{Objects: ids, Quiet: aws.Bool(true)},
		})
		if err != nil {
			return fmt.Errorf("delete %d object(s) from %s: %w", len(batch), batch[0], err)
		}
		if len(out.Errors) > 0 {
			e := out.Errors[0]
			return fmt.Errorf("delete %s: %s: %s (%d key(s) failed in this batch)", aws.ToString(e.Key), aws.ToString(e.Code), aws.ToString(e.Message), len(out.Errors))
		}
	}
	return nil
}
