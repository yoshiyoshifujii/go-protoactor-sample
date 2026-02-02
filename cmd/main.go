package main

import (
	"fmt"
	"time"

	"yoshiyoshifujii/go-protoactor-sample/internal/domain"
	"yoshiyoshifujii/go-protoactor-sample/internal/interface_adaptor"
	"yoshiyoshifujii/go-protoactor-sample/internal/usecase"

	"github.com/asynkron/protoactor-go/actor"
)

func main() {
	system := actor.NewActorSystem()
	root := system.Root

	pingPongUsecase := usecase.NewPingPongUsecase()
	props := actor.PropsFromProducer(interfaceadaptor.NewPingPongProducer(pingPongUsecase))

	pid := root.Spawn(props)
	defer root.Stop(pid)

	future := root.RequestFuture(pid, domain.Ping{Value: "ping"}, time.Second)
	result, err := future.Result()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	pong, ok := result.(domain.Pong)
	if !ok {
		fmt.Printf("unexpected response: %T\n", result)
		return
	}

	fmt.Printf("received: %s\n", pong.Value)
}
