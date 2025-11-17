package domain

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/caddyserver/caddy"
	"github.com/gorilla/websocket"
)

func init() {
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(loadConfig))
}

func loadConfig(serverType string) (caddy.Input, error) {
	contents, err := ioutil.ReadFile(caddy.DefaultConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	fmt.Printf("Loading Caddyfile: %s\n", string(contents))
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       caddy.DefaultConfigFile,
		ServerTypeName: serverType,
	}, nil
}

type Server struct {
	Upgrader websocket.Upgrader
	terminal *Terminal
	Port     int
	Domain   string
}

func NewServer(terminal *Terminal, port int, domain string) *Server {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return &Server{
		Upgrader: upgrader,
		terminal: terminal,
		Port:     port,
		Domain:   domain,
	}
}

func (s *Server) Start() {
	http.HandleFunc("/wssss00", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.Upgrader.Upgrade(w, r, nil) // Upgrade to WebSocket connection
		if err != nil {
			log.Println("Upgrade Error:", err)
			return
		}
		defer conn.Close()

		log.Println("Client Connected")

		// Continuously read messages from the client
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read Error:", err)
				break
			}

			log.Printf("Received Message: %s", message)

			// Echo the message back to the client
			if err := conn.WriteMessage(messageType, message); err != nil {
				log.Println("Write Error:", err)
				break
			}
		}
	})

	serverAddr := fmt.Sprintf("%s:%d", "0.0.0.0", s.Port)
	log.Println("WebSocket server started on", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("ListenAndServe Error:", err)
	}
}

func (s *Server) CaddyServer() {
	caddy.AppName = "Terminal"
	caddy.AppVersion = "0.1"

	caddyfile, err := caddy.LoadCaddyfile("http")
	if err != nil {
		panic(err)
	}

	inst, err := caddy.Start(caddyfile)
	if err != nil {
		panic(err)
	}

	inst.Wait()
}
