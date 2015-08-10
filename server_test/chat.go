package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"sync"
	"text/template"

	"code.google.com/p/go.net/websocket"
)

var (
	httpAddr = flag.String("http", ":8080", "address to serve http")
	apiAddr  = flag.String("api", ":8081", "address to serve api")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", handleWelcome)
	http.HandleFunc("/chat", handleChat)
	go http.ListenAndServe(*httpAddr, nil)

	http.Handle("/ws", websocket.Handler(hendleWebsocket))

	ln, err := net.Listen("tcp", *apiAddr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				log.Fatal(err)
			}
			// go http
			log.Fatal(c)
		}
	}()
}

type Message struct {
	From string
	Text string
}

type Client interface {
	Send(m Message)
}

var (
	clients      = make(map[Client]bool)
	clientsMutex sync.Mutex
)

func registerClient(c Client) {
	clientsMutex.Lock()
	clients[c] = true
	clientsMutex.Unlock()
}

func unregisterClient(c Client) {
	clientsMutex.Lock()
	delete(clients, c)
	clientsMutex.Unlock()
}

func broadcastMessage(m Message) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for c := range clients {
		c.Send(m)
	}

}

var indexPageTempl = template.Must(template.New("").Parse(`
		<html>
			<body>
				<b> Welcome to chat!</b>
				<form action='/chat'>
					<input type="text" id="name" name="name"></input>
					<input type="button" value="Say" onclick="sendMessage()"></input>
				</form> 
				<textarea readonly=1 rows=20 id="alltext"></textarea>
			</body>
		</html>
	`))

var chatPageTempl = template.Must(template.New("").Parse(`
		<html>
			<body>
				<b> Hi,  {{.Name}}!</b>
				<form>
					<input type="text" id="chattext"></input>
					<input type="button" value="Say" onclick="sendMessage()"></input>
				</form> 
				<textarea readonly=1 rows=20 id="alltext"></textarea>
			</body>
		</html>
	`))

func handleWelcome(w http.ResponseWriter, r *http.Request) {

	indexPageTempl.Execute(w, nil)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	log.Printf("serving chat pare for '%v'", name)
	type Params struct {
		Name string
	}

	chatPageTempl.Execute(w, &Params{name})
}

type WSClient struct {
	conn *websocket.Conn
	enc  *json.Encoder
}

func (c *WSClient) Send(m Message) {
	c.enc.Encode(m)
}

func hendleWebsocket(ws *websocket.Conn) {
	c := &WSClient{ws, json.NewEncoder(ws)}
	registerClient(c)
	dec := json.NewDecoder(ws)
	for {
		var m Message
		if err := dec.Decode(&m); err != nil {
			log.Printf("error reading from websocket: %v", err)
			break
		}
		broadcastMessage(m)
	}
	unregisterClient(c)
}
