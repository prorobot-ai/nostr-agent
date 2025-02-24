package bot

import (
	"agent/core"
	"sync"
)

type EventBus struct {
	subscribers map[core.EventType][]func(*core.OutgoingMessage)
	lock        sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[core.EventType][]func(*core.OutgoingMessage)),
	}
}

func (bus *EventBus) Subscribe(eventType core.EventType, handler func(*core.OutgoingMessage)) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	bus.subscribers[eventType] = append(bus.subscribers[eventType], handler)
}

func (bus *EventBus) Publish(eventType core.EventType, message *core.OutgoingMessage) {
	bus.lock.RLock()
	defer bus.lock.RUnlock()

	if handlers, found := bus.subscribers[eventType]; found {
		for _, handler := range handlers {
			go handler(message) // Asynchronous execution
		}
	}
}
