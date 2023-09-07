package main

import (
	"net/rpc"
	"strconv"
	"time"
)

func main() {
	servers := make([]Server, 3)
	for i := 0; i < 3; i++ {
		servers[i].Ports = make([]int, 3)
		servers[i].Myport = 1234 + i
		servers[i].Name = strconv.Itoa(i)
		rpc.RegisterName("Server_"+strconv.Itoa(i), &servers[i]) // need to register each individual object
	}
	for i := 0; i < 3; i++ {

		for j := 0; j < 3; j++ {
			servers[j].Ports[i] = 1234 + i
		}
	}
	for i := 0; i < 3; i++ {
		go servers[i].Listen()
	}
	time.Sleep(2 * time.Second)
	servers[0].Broadcast("hello")
	servers[1].Broadcast("goodbye")
	servers[2].Broadcast("hello again!")
	// Register the Calculator object with the RPC server.

	time.Sleep(3 * time.Second)

}
