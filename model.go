package socket

type WebSocketMsg struct {
	ID       string
	Messages []Message
}

type Message struct {
	Log string
}

type Pipeline struct {
	ID      string `json:"id"`
	Type    int    `json:"type" binding:"required"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var ActivePipelines = make(chan *Pipeline, 6)
var AllPipelines = []*Pipeline{}
