package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]Peer) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var avatarID = 0

type Peer struct {
	Room     string `json:"room"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// Message object
type Message struct {
	Room     string `json:"room"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Message  string `json:"message"`
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "my 404 page!")
}

func FileServerWithCustom404(fs http.FileSystem) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("URL %s\n", r.URL.String())

		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			NotFound(w, r)
			//r.URL.Path = "/" //Serve root
			//fsh.ServeHTTP(w, r)
			return
		}
		fsh.ServeHTTP(w, r)
	})
}

func main() {
	port := 80
	if len(os.Args) > 1 {
		var err error
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			port = 80
		}
		if port == 0 {
			port = 80
		}
	}
	// Use binary asset FileServer
	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "public"}))

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	log.Printf("http server started on :%d\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error %v", err)
		return
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = Peer{"", "", ""}

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		if clients[ws].Username == "" {
			peer := Peer{msg.Room, msg.Username, strconv.Itoa((avatarID%5)+1) + ".png"}
			//peer = Peer(msg.Room, msg.Username, strconv.Itoa((avatarID%5)+1))
			clients[ws] = peer
			//log.Printf("user %s, id %d = %s\n", msg.Username, avatarID, peer.Avatar)
			avatarID = avatarID + 1
		}
		//Don't forward empty msg, but client must send empty msgs to pass Room&User 2022/11
		if msg.Message != "" {
			msg.Avatar = clients[ws].Avatar
			// Send the newly received message to the broadcast channel
			broadcast <- msg
		}
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			if msg.Room == clients[client].Room { //only send to client in same room 2022/11/22
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
