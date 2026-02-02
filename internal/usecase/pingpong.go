package usecase

import "yoshiyoshifujii/go-protoactor-sample/internal/domain"

type PingPongUsecase struct{}

func NewPingPongUsecase() *PingPongUsecase {
	return &PingPongUsecase{}
}

func (u *PingPongUsecase) Handle(ping domain.Ping) domain.Pong {
	return domain.Pong{Value: "pong"}
}
