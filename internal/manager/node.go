package manager

import (
	"log"
	"net"

	"github.com/erkrnt/symphony/api"
	"github.com/erkrnt/symphony/internal/pkg/cluster"
	"github.com/erkrnt/symphony/internal/pkg/config"
	"google.golang.org/grpc"
)

// Node : manager node
type Node struct {
	Key   *config.Key
	Raft  *cluster.RaftNode
	State *cluster.RaftState
}

// NewNode : creates a new manager struct
func NewNode(flags *config.Flags) (*Node, error) {
	k, err := config.GetKey(flags.ConfigDir)

	if err != nil {
		return nil, err
	}

	node, store := cluster.NewRaft(flags)

	m := &Node{
		Key:   k,
		Raft:  node,
		State: store,
	}

	return m, nil
}

// Start : starts Raft memebership server
func Start(f *config.Flags, n *Node) {
	lis, err := net.Listen("tcp", f.ListenRemoteAPI.String())

	if err != nil {
		log.Fatal("Failed to listen")
	}

	s := grpc.NewServer()

	server := &managerServer{
		Node: n,
	}

	api.RegisterManagerServer(s, server)

	log.Print("Started manager gRPC endpoints.")

	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve")
	}
}
