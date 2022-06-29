package socket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func homePage(c *gin.Context) {
	//Upgrade get request to webSocket protocol
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("error get connection")
		log.Fatal(err)
	}
	log.Println(ws)
	// defer ws.Close()
}

func Serversetup() {
	fmt.Println("Hello World")

	r := gin.Default()
	r.GET("/", homePage)
	// r.GET("/ws", wsEndpoint)

	if err := r.Run(":8081"); err != nil {
		log.Fatal(err.Error())
	}
	// setupRoutes()
	// log.Fatal(http.ListenAndServe(":8081", nil))
}

//Runs as an independent goroutine to process pipelines
func StartEngine(pipelines <-chan *Pipeline) {
	fmt.Println("\n I am insdie StartEngine func")
	fmt.Println("\n pipelines len", len(pipelines))
	for p := range pipelines {
		processPipeline(p)
	}

}

func processPipeline(p *Pipeline) {
	fmt.Println("starting to process pipeline having following id ", p.ID, p)
	MarshallingData(p.ID, "starting to process pipeline having following id ")
}

func CreatePipeline(context *gin.Context) {
	fmt.Println("\n I am inside createPipelines")
	var pipelineMsg Pipeline
	if err := context.ShouldBindJSON(&pipelineMsg); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//generate an id and fill it int the object.
	pipelineID := xid.New()
	// log.Info("adding new pipeline to queue with ID ", pipelineID)
	pipelineMsg.ID = pipelineID.String()
	fmt.Println("\n pipelineMsg", pipelineMsg)

	//add pipeline to the queue for execution.
	//while adding also check if the queue is full.
	select {
	case ActivePipelines <- &pipelineMsg:
		fmt.Println("\n I m inside Activepipelines case")
		fmt.Println("\n ActivePipelines len", len(ActivePipelines))
		//store pipeline details in arrawy
		AllPipelines = append(AllPipelines, &pipelineMsg)
		fmt.Println("\n AllPipelines len", len(AllPipelines))
		//send the response back.
		context.IndentedJSON(http.StatusCreated, pipelineMsg)
	default:
		context.IndentedJSON(http.StatusInternalServerError, "Max pipeline limit reached")
	}
}

func MarshallingData(id, logData string) { // []byte {
	fmt.Println("\n I am inside MarshallingData func")
	socketUrl := "ws://localhost:8081" + "/"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
		fmt.Println("Error connecting to Websocket Server:", err)
	}
	var msg Message
	var ws WebSocketMsg
	msg.Log = logData
	ws.ID = id
	ws.Messages = []Message{msg}
	marshalled, _ := json.Marshal(ws)
	// return marshalled

	fmt.Println("\n ws", ws)

	// This is printing at client side
	if err := conn.WriteMessage(1, marshalled); err != nil {
		log.Println(err)
		fmt.Println("err in writing", err)
		return
	}
}

/*
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// helpful log statement to show connections
	log.Println("Client Connected")

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

*/
