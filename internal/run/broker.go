package run

import (
	"sync"
)

type (
	// Broker is a brokerage for publishers and subscribers of Runs.
	Broker struct {
		subscribers []func(event Run)
		mu          sync.RWMutex
	}

	Callback func(event Run)

	Subscriber interface {
		Subscribe(cb Callback)
	}

	Publisher interface {
		Publish(Run)
	}
)

func (b *Broker) Subscribe(cb Callback) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.subscribers = append(b.subscribers, cb)
}

func (b *Broker) Publish(event Run) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, sub := range b.subscribers {
		go sub(event)
	}
}
