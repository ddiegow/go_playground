package main

import (
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)

// Define a struct that will be used for RPC.
type Server struct {
	Ports  []int
	Myport int
	Name   string
}

// Define an RPC method that sends a command
func (c *Server) SendCommand(args *CommandArgs, reply *CommandReply) error {
	fmt.Printf("[Server #%s] received command %s\n", c.Name, args.Command)
	reply.Ok = true
	return nil
}

// Define a struct to represent the arguments of the SendCommand method.
type CommandArgs struct {
	Command string
}
type CommandReply struct {
	Ok bool
}

func (c *Server) Listen() {

	// Create a listener for incoming connections.
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(c.Myport))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("[Server #%s] is listening on port %d...\n", c.Name, c.Myport)

	for {
		// Accept incoming connections.
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle each incoming connection in a separate goroutine.
		go rpc.ServeConn(conn)
	}
}
func (c *Server) SendRPC(rpcName string, args *CommandArgs, reply *CommandReply, address string, snumber int) {

	client, err := rpc.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer client.Close()
	callString := "Server_" + strconv.Itoa(snumber) + "." + rpcName
	err = client.Call(callString, &args, &reply)
	if err != nil {
		fmt.Println("Error calling remote method:", err)
		return
	}
}

func (c *Server) Broadcast(command string) {
	for i := 0; i < 3; i++ { // for each server
		if strconv.Itoa(i) == c.Name {
			continue
		}
		address := "localhost:" + strconv.Itoa(c.Ports[i])
		// Prepare the arguments for the SendCommand method.
		args := CommandArgs{Command: command}
		// Prepare the reply for the SendCommand method
		reply := CommandReply{}
		// Send the rpc
		c.SendRPC("SendCommand", &args, &reply, address, i)
		// Inform user of reply
		fmt.Printf("[Server #%s] received Ok value %v after sending command %s to server #%d\n", c.Name, reply.Ok, command, i)
	}
}
