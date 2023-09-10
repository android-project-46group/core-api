package remote

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func (r *remote) GetImage(ctx context.Context, url string) (io.ReadCloser, error) {
	reader := bytes.NewReader([]byte{})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequestWithContext: %w", err)
	}

	//nolint:bodyclose
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to client.Get: %w", err)
	}

	return resp.Body, nil
}
