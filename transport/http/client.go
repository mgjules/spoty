package http

import (
	"github.com/JulesMike/spoty/config"
	"github.com/JulesMike/spoty/json"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Client is a simple wrapper around resty.Client.
type Client struct {
	*resty.Client
}

// NewClient creates a new Client.
func NewClient(cfg *config.Config) *Client {
	client := resty.New()
	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal
	if !cfg.Prod {
		client.EnableTrace()
		client.SetDebug(true)
	}
	client.OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
		otel.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(r.Header))

		return nil
	})

	return &Client{client}
}
