package dfuse

type Hub struct {
	cli       *wssClient
	broadcast chan []byte

	subscribers map[*subscribe]bool
	register    chan *subscribe
	unregister  chan *subscribe
}

func newHub(cli *wssClient) *Hub {
	return &Hub{
		cli:         cli,
		broadcast:   make(chan []byte),
		subscribers: make(map[*subscribe]bool),

		register:   make(chan *subscribe),
		unregister: make(chan *subscribe),
	}
}

func (h *Hub) run() {
	for {
		select {
		case subscriber := <-h.register:
			h.subscribers[subscriber] = true

		case client := <-h.unregister:
			if _, exist := h.subscribers[client]; exist {
				delete(h.subscribers, client)
			}

		case msg := <-h.broadcast:
			for client, exist := range h.subscribers {
				if exist {
					client.callback(msg)
				} else {
					delete(h.subscribers, client)
				}
			}
		}
	}
}
