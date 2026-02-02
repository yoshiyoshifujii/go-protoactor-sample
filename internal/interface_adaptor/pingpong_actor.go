package interfaceadaptor

import (
	"yoshiyoshifujii/go-protoactor-sample/internal/domain"
	"yoshiyoshifujii/go-protoactor-sample/internal/usecase"

	"github.com/asynkron/protoactor-go/actor"
)

type PingPongActor struct {
	usecase *usecase.PingPongUsecase
}

func NewPingPongActor(usecase *usecase.PingPongUsecase) *PingPongActor {
	return &PingPongActor{usecase: usecase}
}

func NewPingPongProducer(usecase *usecase.PingPongUsecase) actor.Producer {
	return func() actor.Actor {
		return NewPingPongActor(usecase)
	}
}

func (a *PingPongActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case domain.Ping:
		pong := a.usecase.Handle(msg)
		ctx.Respond(pong)
	}
}
