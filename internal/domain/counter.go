package domain

type CounterAdd struct {
	Delta int64
}

type CounterGet struct{}

type CounterValue struct {
	Value int64
}

type Counter struct {
	count int64
}

func NewCounter() Counter {
	return Counter{}
}

func (c *Counter) Add(delta int64) int64 {
	c.count += delta
	return c.count
}

func (c *Counter) ApplyIncrement(delta int64) {
	c.count += delta
}

func (c *Counter) ApplySnapshot(value int64) {
	c.count = value
}

func (c *Counter) Value() int64 {
	return c.count
}
