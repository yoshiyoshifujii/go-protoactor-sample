package interfaceadaptor

import (
	"fmt"
	"math"
	"testing"
	"time"

	"yoshiyoshifujii/go-protoactor-sample/internal/domain"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
	"google.golang.org/protobuf/types/known/anypb"
)

type inMemoryProvider struct {
	state persistence.ProviderState
}

func (p *inMemoryProvider) GetState() persistence.ProviderState {
	return p.state
}

func TestCounterActor_PersistsAndRecovers(t *testing.T) {
	tests := []struct {
		name           string
		deltas         []int64
		expectedBefore int64
		expectedAfter  int64
	}{
		{
			name:           "sum persists and recovers",
			deltas:         []int64{1, 2},
			expectedBefore: 3,
			expectedAfter:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			system := actor.NewActorSystem()
			root := system.Root
			provider := &inMemoryProvider{state: persistence.NewInMemoryProvider(100)}
			props := actor.PropsFromProducer(
				NewCounterProducer(),
				actor.WithReceiverMiddleware(persistence.Using(provider)),
			)

			pid, err := root.SpawnNamed(props, "counter-test")
			if err != nil {
				t.Fatalf("spawn error: %v", err)
			}

			// When
			for _, delta := range tt.deltas {
				value, err := requestValue(root, pid, domain.CounterAdd{Delta: delta})
				if err != nil {
					t.Fatalf("add error: %v", err)
				}
				_ = value
			}

			value, err := requestValue(root, pid, domain.CounterGet{})
			if err != nil {
				t.Fatalf("get error: %v", err)
			}

			if err := root.PoisonFuture(pid).Wait(); err != nil {
				t.Fatalf("stop error: %v", err)
			}

			pid, err = root.SpawnNamed(props, "counter-test")
			if err != nil {
				t.Fatalf("respawn error: %v", err)
			}
			t.Cleanup(func() {
				_ = root.PoisonFuture(pid).Wait()
			})

			recovered, err := requestValue(root, pid, domain.CounterGet{})
			if err != nil {
				t.Fatalf("get error: %v", err)
			}

			// Then
			if value.Value != tt.expectedBefore {
				t.Fatalf("expected %d before recovery, got %d", tt.expectedBefore, value.Value)
			}
			if recovered.Value != tt.expectedAfter {
				t.Fatalf("expected %d after recovery, got %d", tt.expectedAfter, recovered.Value)
			}
		})
	}
}

func TestCounterActor_ApplyPersisted(t *testing.T) {
	tests := []struct {
		name          string
		messages      []*anypb.Any
		expectedCount int64
	}{
		{
			name:          "event, snapshot, event",
			messages:      []*anypb.Any{newCounterEvent(2), newCounterSnapshot(10), newCounterEvent(3)},
			expectedCount: 13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			actor := NewCounterActor()

			// When
			for _, msg := range tt.messages {
				actor.applyPersisted(msg)
			}

			// Then
			if actor.counter.Value() != tt.expectedCount {
				t.Fatalf("expected %d, got %d", tt.expectedCount, actor.counter.Value())
			}
		})
	}
}

func TestCounterActor_ApplyPersisted_InvalidData(t *testing.T) {
	tests := []struct {
		name     string
		message  *anypb.Any
		expected int64
	}{
		{
			name: "invalid varint is ignored",
			message: &anypb.Any{
				TypeUrl: counterEventIncremented,
				Value:   []byte{0x80},
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			actor := NewCounterActor()

			// When
			actor.applyPersisted(tt.message)

			// Then
			if actor.counter.Value() != tt.expected {
				t.Fatalf("expected count to remain %d, got %d", tt.expected, actor.counter.Value())
			}
		})
	}
}

func TestEncodeDecodeInt64RoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		value int64
	}{
		{name: "zero", value: 0},
		{name: "one", value: 1},
		{name: "minus one", value: -1},
		{name: "positive", value: 42},
		{name: "max", value: math.MaxInt64},
		{name: "min", value: math.MinInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			value := tt.value

			// When
			encoded := encodeInt64(value)
			decoded, err := decodeInt64(encoded)

			// Then
			if err != nil {
				t.Fatalf("decode error for %d: %v", value, err)
			}
			if decoded != value {
				t.Fatalf("expected %d, got %d", value, decoded)
			}
		})
	}
}

func TestDecodeInt64Invalid(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{name: "invalid varint", data: []byte{0x80}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			data := tt.data

			// When
			_, err := decodeInt64(data)

			// Then
			if err == nil {
				t.Fatal("expected error for invalid varint")
			}
		})
	}
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
