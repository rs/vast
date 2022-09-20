package vast

import (
	"context"
	"io"
	"net/http"
	"time"
)

type VASTResponse struct {
	Body []byte
	Err  error
}

// VASTResponseHandler handles the repsonse from the request
func ResponseHandler(res *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// LoadVAST attempts to load VAST
func LoadURI(VASTAdTagURI string, milliseconds time.Duration) (*VAST, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*milliseconds)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, VASTAdTagURI, nil)
	if err != nil {
		return nil, err
	}
	body, err := ResponseHandler(http.DefaultClient.Do(req))
	if err != nil {
		return nil, err
	}
	vast, err := UnmarshalVAST(string(body))
	if err != nil {
		return nil, err
	}
	return vast, nil
}
