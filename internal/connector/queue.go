package connector

import (
	"sync"

	"lol-teammate-helper/internal/types"
)

// eventQueue holds at most one pending message per URI, preserving the order
// in which URIs first became pending. A newer message replaces a stale one.
type eventQueue struct {
	mu      sync.Mutex
	pending map[string]types.WSMessageType
	order   []string
	closed  bool
	wake    chan struct{}
}

func newEventQueue() *eventQueue {
	return &eventQueue{
		pending: make(map[string]types.WSMessageType),
		wake:    make(chan struct{}, 1),
	}
}

func (q *eventQueue) push(msg types.WSMessageType) {
	q.mu.Lock()
	if _, queued := q.pending[msg.Uri]; !queued {
		q.order = append(q.order, msg.Uri)
	}
	q.pending[msg.Uri] = msg
	q.mu.Unlock()

	select {
	case q.wake <- struct{}{}:
	default:
	}
}

func (q *eventQueue) close() {
	q.mu.Lock()
	q.closed = true
	q.mu.Unlock()
	select {
	case q.wake <- struct{}{}:
	default:
	}
}

// pop returns the oldest pending message. ok is false when nothing is pending.
func (q *eventQueue) pop() (types.WSMessageType, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.order) == 0 {
		return types.WSMessageType{}, false
	}
	uri := q.order[0]
	q.order = q.order[1:]
	msg := q.pending[uri]
	delete(q.pending, uri)
	return msg, true
}

// out streams messages until the queue is closed and drained.
func (q *eventQueue) out() <-chan types.WSMessageType {
	ch := make(chan types.WSMessageType)
	go func() {
		defer close(ch)
		for {
			if msg, ok := q.pop(); ok {
				ch <- msg
				continue
			}
			q.mu.Lock()
			closed := q.closed
			q.mu.Unlock()
			if closed {
				return
			}
			<-q.wake
		}
	}()
	return ch
}
