package main

import (
	"net"
	"github.com/Sirupsen/logrus"
	"log"
	"os"
)

type Server struct {
	Address string
	Port string
}

func NewServer(addr,port string)*Server{
	server:=&Server{
		Address:addr,
		Port:port,
	}
	return server
}

func (s *Server)Run(){
	logrus.Info("Listening for connections on: ",s.Address+":"+s.Port)
	listener,err:=net.Listen("tcp",s.Address+":"+s.Port)
	if err!=nil{
		logrus.Fatal("Error creating listener: ",err)
	}
	defer listener.Close()
	for{
		conn,err:=listener.Accept()
		if err!=nil{
			logrus.Error("Error accepting connection: ",err)
		}
		c:=&Client{
			Connection:conn,
		}
		go c.HandleConnection()
	}
}

type Client struct {
	Connection net.Conn
}

//GetData read data from connection
func (c *Client)GetData()string{
	data:=make([]byte,1024)
	n,err:=c.Connection.Read(data)
	if err!=nil{
		log.Fatal("Error reading data from client: ",err)
	}
	return string(data[:n])
}

func (c *Client)HandleConnection(){
	request:=c.GetData()
	logrus.Info(request)

	file,err:=os.Create("request.txt")
	if err!=nil{
		logrus.Error("Error creating file: ",err)
	}
	file.WriteString(request)
	file.Close()
}


func main(){
	tcp:=NewServer("","8080")
	tcp.Run()
}
