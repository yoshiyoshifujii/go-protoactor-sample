package interfaceadaptor

import (
	"encoding/binary"
	"fmt"

	"yoshiyoshifujii/go-protoactor-sample/internal/domain"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	counterEventIncremented = "sample.counter/event/Incremented"
	counterSnapshotState    = "sample.counter/snapshot/State"
)

type CounterActor struct {
	persistence.Mixin
	counter domain.Counter
}

func NewCounterActor() *CounterActor {
	return &CounterActor{
		counter: domain.NewCounter(),
	}
}

func (a *CounterActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *persistence.RequestSnapshot:
		a.PersistSnapshot(newCounterSnapshot(a.counter.Value()))
	case *anypb.Any:
		a.applyPersisted(msg)
	case domain.CounterAdd:
		sender := ctx.Sender()
		if !a.Recovering() {
			a.PersistReceive(newCounterEvent(msg.Delta))
		}
		current := a.counter.Add(msg.Delta)
		if sender != nil {
			ctx.Send(sender, domain.CounterValue{Value: current})
		}
	case domain.CounterGet:
		ctx.Respond(domain.CounterValue{Value: a.counter.Value()})
	}
}

func (a *CounterActor) applyPersisted(msg *anypb.Any) {
	value, err := decodeInt64(msg.Value)
	if err != nil {
		return
	}
	switch msg.TypeUrl {
	case counterSnapshotState:
		a.counter.ApplySnapshot(value)
	case counterEventIncremented:
		a.counter.ApplyIncrement(value)
	}
}

func NewCounterProducer() func() actor.Actor {
	return func() actor.Actor {
		return NewCounterActor()
	}
}

func newCounterEvent(delta int64) *anypb.Any {
	return &anypb.Any{
		TypeUrl: counterEventIncremented,
		Value:   encodeInt64(delta),
	}
}

func newCounterSnapshot(count int64) *anypb.Any {
	return &anypb.Any{
		TypeUrl: counterSnapshotState,
		Value:   encodeInt64(count),
	}
}

func encodeInt64(v int64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, v)
	return buf[:n]
}

func decodeInt64(data []byte) (int64, error) {
	value, n := binary.Varint(data)
	if n <= 0 {
		return 0, fmt.Errorf("invalid varint")
	}
	return value, nil
}
