package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Client struct {
	Connection net.Conn
}

func (c *Client) Close() error {
	return c.Connection.Close()
}

// ReadBody read body from connection
func (c *Client) ReadBody() (string, error) {
	data := make([]byte, 1024)
	n, err := c.Connection.Read(data)
	if err != nil || n == 0 {
		return "", err
	}
	return string(data[:n]), err
}

func (c *Client) WriteBytes(msg []byte) error {
	_, err := c.Connection.Write(msg)
	return err
}

func (c *Client) WriteString(msg string) error {
	_, err := c.Connection.Write([]byte(msg))
	return err
}

func (c *Client) HandleConnection() error {
	requestBody, err := c.ReadBody()
	if err != nil {
		return err
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Println("cannot close connection: ", err)
		}
	}()

	//creates temporary file to store request
	f, err := os.Create("request.txt")
	if err != nil {
		return fmt.Errorf("couldn't create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	if _, err = f.WriteString(requestBody); err != nil {
		return fmt.Errorf("couldn't write request to file: %v", err)
	}

	req, err := convertToRequest(f.Name())
	if err != nil {
		return err
	}

	res, err := createResponse(req)
	if err != nil {
		return err
	}

	c.WriteBytes(res.toByte())

	return nil
}
