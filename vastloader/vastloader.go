package vastloader

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/KargoGlobal/vast"
	vastXml "github.com/KargoGlobal/vast"
)

type VASTResponse struct {
	Body []byte
	Err  error
}

// VASTResponseHandler handles the repsonse from the request
func VASTResponseHandler(res *http.Response, err error) VASTResponse {
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
func LoadVAST(VASTAdTagURI string) (*vastXml.VAST, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, VASTAdTagURI, nil)
	if err != nil {
		return nil, err
	}
	responseChan := make(chan VASTResponse, 0)
	go func() {
		responseChan <- VASTResponseHandler(http.DefaultClient.Do(req))
	}()
	vastResponse := <-responseChan
	if vastResponse.Err != nil {
		return nil, vastResponse.Err
	}
	vast, err := vast.UnmarshalVAST(string(vastResponse.Body))
	if err != nil {
		return nil, err
	}
	return vast, nil
}
