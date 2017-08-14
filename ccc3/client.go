package main

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/clientselector"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

func main() {
	server1 := &clientselector.ServerPeer{Network: "tcp", Address: "127.0.0.1:8972"}
	server2 := &clientselector.ServerPeer{Network: "tcp", Address: "127.0.0.1:8973"}

	servers := []*clientselector.ServerPeer{server1, server2}

	s := clientselector.NewMultiClientSelector(servers, rpcx.RoundRobin, 10*time.Second)
	s.SelectMode = rpcx.RandomSelect
	client := rpcx.NewClient(s)
	client.Timeout = 6 * time.Second

	// compressionPlugin := plugin.NewCompressionPlugin(rpcx.CompressSnappy)
	// client.PluginContainer.Add(compressionPlugin)

	//rpcx.Reconnect = nil

	for i := 0; i < 10; i++ {
		callServer(client)
		time.Sleep(6 * time.Second)
	}

	client.Close()
}

func callServer(client *rpcx.Client) {
	args := &Args{7, 8}
	var reply Reply
	err := client.Call(context.Background(), "Arith.Mul", args, &reply)
	if err != nil {
		log.Infof("error for Arith: %d*%d, %v", args.A, args.B, err)
	} else {
		log.Infof("Arith: %d*%d=%d", args.A, args.B, reply.C)
	}
}
