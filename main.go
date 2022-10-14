package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	Address string
	Port    string
}

func NewServer(addr, port string) *Server {
	return &Server{
		Address: addr,
		Port:    port,
	}
}

func (s *Server) Run() {
	log.Printf("Server working on: %s:%s", s.Address, s.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.Address, s.Port))
	if err != nil {
		log.Fatal("Error creating listener: ", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err)
		}
		c := &Client{
			Connection: conn,
		}
		go func() {
			if err := c.HandleConnection(); err != nil {
				log.Println(err)
			}
		}()
	}
}

func main() {
	tcp := NewServer("", "8080")
	tcp.Run()
}
