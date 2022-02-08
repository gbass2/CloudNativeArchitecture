// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
	"strings"
	"os"
)

// Client struct which holds the client's name and channel
type Client struct {
    name string
		conn net.Conn
    cli chan<- string
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

type Clients struct {
	clients map[Client]bool // all connected clients
	chatHistory []string
}

// Default constructor for Clients
func newClients() *Clients {
	c:= new(Clients)
	c.clients = make(map[Client]bool)
	return c
}
// Prints the online clients
func (c *Clients) printOnlineUsers(cli Client) {
	// Print the members that are online to the channel
	cli.cli <- "\nMembers Online:"
	for currClient,_ := range(clients.clients){
		cli.cli <- currClient.name
	}

	cli.cli <- ""
}

// Prints the chat history
func (c *Clients) printChatHistory(cli Client) {
	// Print the members that are online to the channel
	cli.cli <- "\nChat History:"
	for _,message := range(clients.chatHistory){
		cli.cli <- message
	}

	cli.cli <- ""
}

// Saves the chat history upon program exit.
func (c *Clients) saveChatHistory() error {
	f, err := os.Create("chat_history.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	for _, message := range clients.chatHistory {
	 fmt.Fprintln(f, message)
	}
	return nil
}

// Direct messages a connected client
func (c *Clients) directMessage(input string, cli Client) {
	// Print the members that are online to the channel
	words := strings.Fields(input)
	name := strings.Replace(words[0], "@", "", -1)
	message := "(" + time.Now().Format("01-02-2006 15:04:05") + ") (DM) " + cli.name + ": "+ input
	debug := "(" + time.Now().Format("01-02-2006 15:04:05") + ") (DM) " + cli.conn.RemoteAddr().String() + ": " + cli.name + ": " + input

	for k,_ := range(clients.clients){
		if name == k.name {
			fmt.Println(debug)
			k.cli <- message
			cli.cli <- message

		}
	}
}


func (c *Clients) printHelp(cli Client) {
	cli.cli <- "\nThe commands are:\n /help - List the commands. \n /online - List the members that are in the chat room.\n @(Person's name) - Use the @ symbol to dm someone. \n"
}

func (c *Clients) commands(input string, cli Client) bool {
	command := true
	switch {
	case input == "/online":
		c.printOnlineUsers(cli)
	case input == "/help":
		c.printHelp(cli)
	case input[0:1] == "@":
		c.directMessage(input, cli)

	default:
		command = false
	}

	return command
}

var (
	entering = make(chan Client)
	leaving  = make(chan Client)
	messages = make(chan string) // all incoming client messages
	clients = newClients()
)

func main() {
	listener, err := net.Listen("tcp", "192.168.1.100:6666")
	if err != nil {
		log.Fatal(err)
	}

	// Saving chat history upon exit
	go func(){
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			if input.Text() == "exit" {
				clients.saveChatHistory()
				os.Exit(1)
			}
		}
	}()

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

func broadcaster() {
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients.clients {
				cli.cli <- msg
			}

		case cli := <-entering:
			clients.clients[cli] = true

			// Print the members that are online to the channel
			clients.printOnlineUsers(cli)

			// Print the chat history to the new client
			if len(clients.chatHistory) > 0 {
				clients.printChatHistory(cli)
			}


		case cli := <-leaving:
			delete(clients.clients, cli)
			close(cli.cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)
	currClient := Client{name: "", conn: conn, cli: ch}

	// Getting the name client's name
	currClient.readName(conn)

	ch <- "Welcome " + currClient.name + ". To get started type '/help' to see the list of commands."

	messages <- "(" + time.Now().Format("01-02-2006 15:04:05") + ") " + currClient.name + " has arrived"
	fmt.Println("(" + time.Now().Format("01-02-2006 15:04:05") + ") " + conn.RemoteAddr().String()+ ": "+ currClient.name + " has arrived")

	entering <- currClient

	input := bufio.NewScanner(conn)
	for input.Scan() {
		command := clients.commands(input.Text(), currClient)
		message := "(" + time.Now().Format("01-02-2006 15:04:05") + ") " + currClient.name + ": " + input.Text()
		if strings.TrimSpace(input.Text()) != "" && !command {
			fmt.Println("(" + time.Now().Format("01-02-2006 15:04:05") + ") " + conn.RemoteAddr().String()+ ": "+ currClient.name + ": " + input.Text())
			messages <- message
			clients.chatHistory = append(clients.chatHistory, message)
		}
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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
