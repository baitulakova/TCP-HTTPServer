package main

import (
	"net"
	"github.com/Sirupsen/logrus"
	"os"
	"github.com/baitulakova/TCP-HTTPServer/request"
)

//temporary file to store request
const filename  ="request.txt"

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
	logrus.Info("Server working on: ",s.Address+":"+s.Port)
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

func (c *Client)Close(){
	c.Connection.Close()
	logrus.Info("Connection with %v closed",c.Connection.RemoteAddr().String())
}

//GetData read data from connection
func (c *Client)GetData()(string,error){
	data:=make([]byte,1024)
	n,err:=c.Connection.Read(data)
	if err!=nil{
		return "",err
	}
	return string(data[:n]),err
}

func (c *Client)HandleConnection(){
	logrus.Info("Serving ",c.Connection.RemoteAddr().String())
	Request,err:=c.GetData()
	if err!=nil{
		c.Close()
	}

	//creates temporary file to store request
	file,err:=os.Create(filename)
	if err!=nil{
		logrus.Error("Error creating file: ",err)
		c.Close()
	}
	_,err=file.WriteString(Request)
	if err!=nil{
		logrus.Error("Error writing request to file")
		c.Close()
	}
	file.Close()
	requestLines:=request.HandleRequest(filename)
	req:=request.FormRequest(requestLines)
	logrus.Info(req)
}

func main(){
	tcp:=NewServer("","8080")
	tcp.Run()
}
