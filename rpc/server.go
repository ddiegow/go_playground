package main

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

// Define a struct that will be used for RPC.
type Server int

// Define an RPC method that adds two numbers.
func (c *Server) SendCommand(args *CommandArgs, reply *CommandReply) error {
	fmt.Printf("Server received command %s\n", args.Command)
	reply.Ok = true
	return nil
}

// Define a struct to represent the arguments of the Add method.
type CommandArgs struct {
	Command string
}
type CommandReply struct {
	Ok bool
}

func (c *Server) Listen() {

	// Create a listener for incoming connections.
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 1234...")

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
func (c *Server) SendRPC(rpcName string) {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer client.Close()

	// Prepare the arguments for the Add method.
	args := CommandArgs{Command: "BecomeLaggy"}

	// Call the remote Add method and store the result.
	reply := CommandReply{}
	callString := "Server." + rpcName
	err = client.Call(callString, &args, &reply)
	if err != nil {
		fmt.Println("Error calling remote method:", err)
		return
	}
	fmt.Printf("Ok value: %v\n", reply.Ok)
}
func main() {
	server := new(Server)

	// Register the Calculator object with the RPC server.
	rpc.Register(server)
	go server.Listen()
	time.Sleep(2 * time.Second)
	go server.SendRPC("SendCommand")
	go server.SendRPC("SendCommand")
	go server.SendRPC("SendCommand")
	time.Sleep(3 * time.Second)

}
