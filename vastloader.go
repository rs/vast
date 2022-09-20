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
func ResponseHandler(res *http.Response, err error) VASTResponse {
	if err != nil {
		return VASTResponse{
			Body: []byte(""),
			Err:  err,
		}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return VASTResponse{
			Body: []byte(""),
			Err:  err,
		}
	}
	return VASTResponse{
		Body: body,
		Err:  nil,
	}
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
	vastResponse := ResponseHandler(http.DefaultClient.Do(req))
	if vastResponse.Err != nil {
		return nil, vastResponse.Err
	}
	vast, err := UnmarshalVAST(string(vastResponse.Body))
	if err != nil {
		return nil, err
	}
	return vast, nil
}
