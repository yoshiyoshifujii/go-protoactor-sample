package domain

type CounterAdd struct {
	Delta int64
}

type CounterGet struct{}

type CounterValue struct {
	Value int64
}
