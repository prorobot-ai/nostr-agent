package bot

import (
	"agent/core"
	"sync"
)

type EventBus struct {
	subscribers map[core.EventType][]func(*core.Message)
	lock        sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[core.EventType][]func(*core.Message)),
	}
}

func (bus *EventBus) Subscribe(eventType core.EventType, handler func(*core.Message)) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)
}

func (bus *EventBus) Publish(eventType core.EventType, message *core.Message) {
	bus.lock.RLock()
	defer bus.lock.RUnlock()

	if handlers, found := bus.subscribers[eventType]; found {
		for _, handler := range handlers {
			go handler(message) // Asynchronous execution
		}
	}
}
