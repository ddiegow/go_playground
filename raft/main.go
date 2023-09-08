package main

import (
	"net/rpc"
	"strconv"
	"time"
)

func main() {
	servers := make([]Server, 3)
	for i := 0; i < 3; i++ {
		servers[i].ports = make([]int, 3)
		servers[i].myPort = 1234 + i
		servers[i].myId = strconv.Itoa(i)
		rpc.RegisterName("Server_"+strconv.Itoa(i), &servers[i]) // need to register each individual object
	}
	for i := 0; i < 3; i++ {

		for j := 0; j < 3; j++ {
			servers[j].ports[i] = 1234 + i
		}
	}
	for i := 0; i < 3; i++ {
		go servers[i].listen()
	}
	time.Sleep(2 * time.Second)

	// Register the Calculator object with the RPC server.

	time.Sleep(3 * time.Second)

}
