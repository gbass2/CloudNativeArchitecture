// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

// Brian explain project
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
	"strings"
)

type client chan<- string // an outgoing message channel

var (
	entering = make(chan Client)
	leaving  = make(chan Client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	listener, err := net.Listen("tcp", "192.168.1.223:9999")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// Noah
func broadcaster() {
	clients := make(map[Client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

			// Print the members that are online to the channel
			cli.cli <- "\nMembers Online:"
			for currClient,_ := range(clients){
				cli.cli <- currClient.name
			}

			cli.cli <- ""

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.cli)
		}
	}
}

// Grayson
func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	currClient := Client{name: "", cli: ch}

	// Getting the name client's name
	currClient.readName(conn)

	ch <- "You are " + currClient.name

	messages <- "(" + time.Now().Format("01-02-2006 15:04:05") + ") " + currClient.name + " has arrived"
	fmt.Println("(" + time.Now().Format("01-02-2006 15:04:05") + ") " + conn.RemoteAddr().String()+ ": "+ currClient.name + " has arrived")

	entering <- currClient

	input := bufio.NewScanner(conn)
	for input.Scan() {
		fmt.Println("(" + time.Now().Format("01-02-2006 15:04:05") + ") " + conn.RemoteAddr().String()+ ": "+ currClient.name + ": " + input.Text())
		messages <- "(" + time.Now().Format("01-02-2006 15:04:05") + ") " + currClient.name + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- currClient
	fmt.Println("(" + time.Now().Format("01-02-2006 15:04:05") + ") " + conn.RemoteAddr().String()+ ": " + currClient.name + " has left")
	messages <- "(" + time.Now().Format("01-02-2006 15:04:05") + ") " + currClient.name + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

// Garrett
// Client struct which holds the client's name and channel
type Client struct {
    name string
    cli client
}

// Returns an entered name for the client
func (c *Client) readName(conn net.Conn) {
	c.cli <- "Enter your name: "
	input := bufio.NewScanner(conn)
	for input.Scan() {
		c.name = input.Text() // Get name from command line
		// If the name is blank or contains only whitespaces then let client re-enter their name
		if strings.TrimSpace(c.name) != "" {
			break
		}
	}
}
