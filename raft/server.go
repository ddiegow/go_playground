package main

import (
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

const (
	FOLLOWER = iota
	CANDIDATE
	LEADER
)

// Define a struct that will be used for RPC.
type Server struct {
	ports  []int
	myPort int
	myId   string

	currentTerm int     // latest term server has seen (initialized to 0 on first boot, increases monotonically)
	votedFor    int     // candidateId that received vote in current term (or null if none)
	log         []Entry // log entries; each entry contains command for state machine, and term when entry was received by leader (first index is 1)
	// Volatile state on all servers:
	commitIndex int // index of highest log entry known to be committed (initialized to 0, increases monotonically)
	lastApplied int // index of highest log entry applied to state 	machine (initialized to 0, increases monotonically)
	// Volatile state on leaders: (Reinitialized after election)
	nextIndex  []int // for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	matchIndex []int // for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)

	serverState int
	mu          sync.Mutex

	hearbeatChan    chan bool // channel to receive heartbeats
	votedChan       chan bool // channel to indicate we have just vote
	electionWonChan chan bool // channel to indicate we won the election and need to convert to leader
	stepDownChan    chan bool // channel to let the leader know it has to step down
}

// we have received a vote request from another server
func (s *Server) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if args.Term < s.currentTerm { //  Reply false if term < currentTerm
		reply.Term = s.currentTerm
		reply.VoteGranted = false
		return nil
	}
	if args.Term > s.currentTerm {
		s.currentTerm = args.Term
		// convert to candidate
	}
	//  If votedFor is null or candidateId, and candidate’s log is at least as up-to-date as receiver’s log, grant vote
	if (s.votedFor == -1 || s.votedFor == args.CandidateId) && s.logUpToDate(args.LastLogTerm, args.LastLogIndex) {
		reply.Term = s.currentTerm
		reply.VoteGranted = true
		return nil
	}
	return nil
}

// Get last log index. Needs to be called with lock in place
func (s *Server) getLastIndex() int {
	return len(s.log) - 1
}

// Get last log term. Needs to be called with lock in place
func (s *Server) getLastTerm() int {
	if s.getLastIndex() > -1 {
		return s.log[s.getLastIndex()].Term
	} else {
		return -1
	}
}

func (s *Server) logUpToDate(lastLogTerm int, lastLogIndex int) bool {
	myLastIndex := s.getLastIndex()
	myLastTerm := s.getLastTerm()
	if lastLogIndex < myLastIndex || myLastTerm > lastLogTerm {
		return false
	}
	return true

}

// This function coordinates the changes from one state to another.
func (s *Server) coordinate() {
	s.mu.Lock()
	currentState := s.serverState
	s.mu.Unlock()
	switch currentState {
	case FOLLOWER:
		select {
		case <-s.votedChan:
		case <-s.hearbeatChan:
		case <-time.After(3 * time.Second): // 3 seconds have passed without a heartbeat
			// convert to candidate
			s.convertToCandidate()
			// start election
			s.startElection()

		}
	case CANDIDATE:
		select {
		case <-s.hearbeatChan: // we got a heartbeat, so convert to follower
		case <-s.electionWonChan: // we won the election, so convert to leader
			s.convertToLeader()
		case <-time.After(3 * time.Second): // 3 seconds have passed without any results
			// start new election
			s.startElection()
		}
	case LEADER:
		select {
		case <-s.stepDownChan: // not the leader anymore, so step down
		// (this is called after stepping down, so next switch iteration it will go into follower block)
		case <-time.After(50 * time.Millisecond):
			// broadcast heartbeats
		}
	}
}

// follower -> candidate transition
func (s *Server) convertToCandidate() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resetChans()
}

// candidate or leader -> follwer transition
func (s *Server) convertToFollower() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resetChans()
}

// candidate -> leader transition
func (s *Server) convertToLeader() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resetChans()
}

// start an election
func (s *Server) startElection() {
}

// restart all channels. must hold lock when accessing this function
func (s *Server) resetChans() {
	s.electionWonChan = make(chan bool)
	s.hearbeatChan = make(chan bool)
	s.stepDownChan = make(chan bool)
	s.votedChan = make(chan bool)
}

// Define an RPC method that sends a command
func (s *Server) sendCommand(args *CommandArgs, reply *CommandReply) error {
	fmt.Printf("[Server #%s] received command %s\n", s.myId, args.Command)
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

func (s *Server) listen() {

	// Create a listener for incoming connections.
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.myPort))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("[Server #%s] is listening on port %d...\n", s.myId, s.myPort)

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
func (s *Server) sendRPC(rpcName string, args *CommandArgs, reply *CommandReply, address string, snumber int) {

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

func (s *Server) broadcast(command string) {
	for i := 0; i < 3; i++ { // for each server
		if strconv.Itoa(i) == s.myId {
			continue
		}
		address := "localhost:" + strconv.Itoa(s.ports[i])
		// Prepare the arguments for the SendCommand method.
		args := CommandArgs{Command: command}
		// Prepare the reply for the SendCommand method
		reply := CommandReply{}
		// Send the rpc
		s.sendRPC("SendCommand", &args, &reply, address, i)
		// Inform user of reply
		fmt.Printf("[Server #%s] received Ok value %v after sending command %s to server #%d\n", s.myId, reply.Ok, command, i)
	}
}
