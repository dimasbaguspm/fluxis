package service

import "github.com/dimasbaguspm/fluxis/pkg/pubsub"

type Deps struct {
	Bus pubsub.Subscriber
}

type Service struct {
	Deps
}

func New(d Deps) *Service {
	return &Service{d}
}
