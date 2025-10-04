package gateway

type Option func(*Gateway) error

func WithBrokerURIs(uris []string) Option {
	return func(gw *Gateway) error {
		gw.brokerURIs = uris
		return nil
	}
}

func ConsumingFromTopic(topic string) Option {
	return func(gw *Gateway) error {
		gw.topic = topic
		return nil
	}
}

func WithLogger(logger Logger) Option {
	return func(gw *Gateway) error {
		gw.logger = logger
		return nil
	}
}

func WithName(name string) Option {
	return func(gw *Gateway) error {
		gw.name = name
		return nil
	}
}

func WithHandler(handler Handler) Option {
	return func(gw *Gateway) error {
		gw.handler = handler
		return nil
	}
}

// func WithExponentialBackoff(initial time.Duration, limit int) Option {
// 	return func(gw *Gateway) error {
// 		gw.backoff = &ExponentialBackoff{
// 			initial: initial,
// 			limit:   limit,
// 		}
// 		return nil
// 	}
// }
