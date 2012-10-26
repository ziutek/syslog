package syslog

// Hanler handles syslog messages
type Handler interface {
	// Handle should return Message (mayby modified) for future processing by
	// other handlers or return nil. If Handle is called with nil message it
	// should complete all remaining work and properly shutdown before return.
	Handle(*Message) *Message
}

// BaseHandler is desigend for simplify the creation of your own handlers. It
// implements Handler interface using nonblocking queuing of messages and
// simple message filtering.
type BaseHandler struct {
	queue  chan *Message
	end    chan struct{}
	filter func(*Message) bool
	ft     bool
}

// NewBaseHandler creates BaseHandler using specified filter. A message is
// queued in BaseHandler if filter is nil or returns true. If ft is true
// message is returned by handler for future processing by other handlers.
func NewBaseHandler(qlen int, filter func(*Message) bool, ft bool) *BaseHandler {
	return &BaseHandler{
		queue:  make(chan *Message, qlen),
		end:    make(chan struct{}),
		filter: filter,
		ft:     ft,
	}
}

// Handle inserts m in an internal queue. If queue is full it immediately
// returns m otherwise it returns nil or m depending on whether ft is false or
// true.
func (h *BaseHandler) Handle(m *Message) *Message {
	if m == nil {
		close(h.queue)
		<-h.end
		return nil
	}
	if h.filter == nil || !h.filter(m) {
		return m
	}
	select {
	case h.queue <- m:
		if !h.ft {
			return nil
		}
	default:
	}
	return m
}

// Get returns first message from internal queue. It waits for message if queue
// is empty. It returns nil if there is no more messages to process.
func (h *BaseHandler) Get() *Message {
	m, ok := <-h.queue
	if ok {
		return m
	}
	return nil
}

// End indicates the server that handler properly shutdown. You need to call End
// only if Get has returned nil before.
func (h *BaseHandler) End() {
	close(h.end)
}
