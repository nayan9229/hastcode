package views

import (
	"context"
	"encoding/json"
	"io"

	"github.com/a-h/templ"
)

// StructuredData writes each JSON-LD block (internal/schema values) as its
// own <script type="application/ld+json">. encoding/json escapes <, > and
// & as < etc., so content can never close the script element.
func StructuredData(blocks []any) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		for _, b := range blocks {
			data, err := json.Marshal(b)
			if err != nil {
				return err
			}
			if _, err := io.WriteString(w, `<script type="application/ld+json">`+string(data)+`</script>`); err != nil {
				return err
			}
		}
		return nil
	})
}
