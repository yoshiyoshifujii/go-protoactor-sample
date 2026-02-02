package main

import (
	"fmt"
	"time"

	"yoshiyoshifujii/go-protoactor-sample/internal/domain"
	"yoshiyoshifujii/go-protoactor-sample/internal/interface_adaptor"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
)

type inMemoryProvider struct {
	state persistence.ProviderState
}

func (p *inMemoryProvider) GetState() persistence.ProviderState {
	return p.state
}

func main() {
	system := actor.NewActorSystem()
	root := system.Root

	provider := &inMemoryProvider{state: persistence.NewInMemoryProvider(3)}
	props := actor.PropsFromProducer(
		interfaceadaptor.NewCounterProducer(),
		actor.WithReceiverMiddleware(persistence.Using(provider)),
	)

	pid, err := root.SpawnNamed(props, "counter")
	if err != nil {
		fmt.Printf("spawn error: %v\n", err)
		return
	}

	addAndPrint(root, pid, 1)
	addAndPrint(root, pid, 2)
	addAndPrint(root, pid, 3)
	getAndPrint(root, pid)

	_ = root.PoisonFuture(pid).Wait()

	pid, err = root.SpawnNamed(props, "counter")
	if err != nil {
		fmt.Printf("respawn error: %v\n", err)
		return
	}

	getAndPrint(root, pid)
	_ = root.PoisonFuture(pid).Wait()
}

func addAndPrint(root *actor.RootContext, pid *actor.PID, delta int64) {
	resp, err := requestValue(root, pid, domain.CounterAdd{Delta: delta})
	if err != nil {
		fmt.Printf("add error: %v\n", err)
		return
	}
	fmt.Printf("add %+d => %d\n", delta, resp.Value)
}

func getAndPrint(root *actor.RootContext, pid *actor.PID) {
	resp, err := requestValue(root, pid, domain.CounterGet{})
	if err != nil {
		fmt.Printf("get error: %v\n", err)
		return
	}
	fmt.Printf("current => %d\n", resp.Value)
}

func requestValue(root *actor.RootContext, pid *actor.PID, msg interface{}) (domain.CounterValue, error) {
	future := root.RequestFuture(pid, msg, time.Second)
	result, err := future.Result()
	if err != nil {
		return domain.CounterValue{}, err
	}
	value, ok := result.(domain.CounterValue)
	if !ok {
		return domain.CounterValue{}, fmt.Errorf("unexpected response: %T", result)
	}
	return value, nil
}
