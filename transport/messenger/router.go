package messenger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/mgjules/spoty/config"
	"github.com/mgjules/spoty/health"
	"github.com/mgjules/spoty/logger"
	"go.uber.org/fx"
)

const (
	retries         = 3
	initialInterval = time.Millisecond * 100
)

// Module exported to initialise a new Publisher, Subscriber and Router.
var Module = fx.Options(
	fx.Provide(NewPublisher, NewSubscriber, NewRouter),
)

// Router is a wrapper for a message router.
type Router struct {
	*message.Router

	publisher  *Publisher
	subscriber *Subscriber

	cfg    *config.Config
	logger *logger.Logger
	health *health.Checks
}

// NewRouter returns a new router.
func NewRouter(
	cfg *config.Config,
	publisher *Publisher,
	subscriber *Subscriber,
	logger *logger.Logger,
	health *health.Checks,
) (*Router, error) {
	r := Router{
		publisher:  publisher,
		subscriber: subscriber,
		cfg:        cfg,
		logger:     logger,
		health:     health,
	}

	wlog := watermill.NewStdLoggerWithOut(logger.Writer(), !cfg.Prod, false)

	router, err := message.NewRouter(message.RouterConfig{}, wlog)
	if err != nil {
		return nil, fmt.Errorf("failed to create message router: %w", err)
	}

	router.AddMiddleware(
		middleware.CorrelationID,

		middleware.Retry{
			MaxRetries:      retries,
			InitialInterval: initialInterval,
			Logger:          wlog,
		}.Middleware,

		middleware.Recoverer,
	)

	r.Router = router
	r.health.RegisterChecks(r.Check())

	return &r, nil
}

// AddHandler is a wrapper around message.Router.AddHandler.
func (r *Router) AddHandler(
	handlerName,
	subscribeTopic,
	publishTopic string,
	handlerFunc message.HandlerFunc,
) *message.Handler {
	return r.Router.AddHandler(
		handlerName,
		subscribeTopic,
		r.subscriber.MessageSubscriber(),
		publishTopic,
		r.publisher.MessagePublisher(),
		handlerFunc,
	)
}

// AddNoPublisherHandler is a wrapper around message.Router.AddNoPublisherHandler.
func (r *Router) AddNoPublisherHandler(
	handlerName,
	subscribeTopic string,
	handlerFunc message.NoPublishHandlerFunc,
) *message.Handler {
	return r.Router.AddNoPublisherHandler(
		handlerName,
		subscribeTopic,
		r.subscriber.MessageSubscriber(),
		handlerFunc,
	)
}

// Publisher returns the publisher for the router.
func (r *Router) Publisher() *Publisher {
	return r.publisher
}

// Subscriber returns the subscriber for the router.
func (r *Router) Subscriber() *Subscriber {
	return r.subscriber
}

// Check is used to perform healthcheck.
func (r *Router) Check() health.Check {
	//nolint:revive
	return health.Check{
		Name:          "messenger.router",
		RefreshPeriod: 10 * time.Second,
		InitialDelay:  10 * time.Second,
		Timeout:       5 * time.Second,
		Check: func(_ context.Context) error {
			if !r.publisher.IsConnected() {
				return errors.New("publisher is not connected")
			} else if !r.subscriber.IsConnected() {
				return errors.New("subscriber is not connected")
			}

			return nil
		},
	}
}
