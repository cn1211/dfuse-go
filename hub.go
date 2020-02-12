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

		register:   make(chan *subscribe, 1),
		unregister: make(chan *subscribe, 1),
	}
}

func (h *Hub) run() {
	for {
		select {
		case subscriber := <-h.register:
			h.subscribers[subscriber] = true

		case subscriber := <-h.unregister:
			if _, exist := h.subscribers[subscriber]; exist {
				delete(h.subscribers, subscriber)
			}

		case msg := <-h.broadcast:
			for client, exist := range h.subscribers {
				if exist {
					client.distribute(msg)
				} else {
					delete(h.subscribers, client)
				}
			}
		}
	}
}
