package ws

// import (
// 	"log"

// 	"github.com/gofiber/websocket/v2"
// )

// // Msg represents the general message structure
// type Msg struct {
// 	Msg  string
// 	Data interface{}
// }

// // Client is a middleman between the websocket connection and the hub.
// type Client struct {
// 	// The websocket connection
// 	Conn *websocket.Conn

// 	// Buffered channel of message
// 	SendMsg chan Msg

// 	// ID of the client (type of client: account, terminal, etc.)
// 	ID string
// }

// // New creates a new instance of Client.
// func New(conn *websocket.Conn) *Client {
// 	return &Client{
// 		Conn:    conn,
// 		SendMsg: make(chan Msg),
// 		ID:      "client",
// 	}
// }

// // Connect starts read and write messages loop.
// func (c *Client) Connect() {
// 	// Envoi des messages
// 	// ------------------
// 	// go c.writeMessage()

// 	// Ecoute des messages
// 	// -------------------
// 	go c.readMessages()
// }

// func (c *Client) readMessages() {
// 	defer func() {
// 		// Déconnexion du hub
// 		// ------------------
// 		err := c.Conn.Close()
// 		log.Printf("Close connection - err=%v\n", err)
// 	}()

// 	// Gestion des messages
// 	// --------------------
// 	// var (
// 	// 	msgType int
// 	// 	msg     []byte
// 	// 	err     error
// 	// )
// 	for {
// 		// if msgType, msg, err = c.Conn.ReadMessage(); err != nil {
// 		// 	log.Printf("[error] read: %v, type=%v, msg=%v", err, msgType, msg)
// 		// 	break
// 		// }

// 		// // Est-ce un JSON valide ?
// 		// // -----------------------
// 		// var msgJSON Msg
// 		// err = json.Unmarshal(bytes.TrimSpace(msg), &msgJSON)

// 		// if err != nil {
// 		// 	// JSON non valide
// 		// 	// ---------------
// 		// 	log.Printf("Read message error: %v\n", err)
// 		// } else {
// 		// 	// JSON valide
// 		// 	// -----------
// 		// 	if msgJSON.Msg != "" {
// 		// 		c.SendMsg <- msgJSON
// 		// 	}
// 		// }
// 	}
// }

// func (c *Client) writeMessage() {
// 	defer func() {
// 		c.Conn.Close()
// 	}()

// 	for {
// 		select {
// 		case msg := <-c.SendMsg:
// 			switch msg.Msg {
// 			case "test":
// 				log.Printf("Write message: %v\n", msg)
// 			}
// 		}
// 	}
// }
