package main

import (
	"context"
	"time"

	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/log"
	"github.com/smallnest/rpcx/plugin"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

func main() {
	s := &rpcx.DirectClientSelector{Network: "tcp", Address: "127.0.0.1:8972", DialTimeout: 10 * time.Second}
	client := rpcx.NewClient(s)
	client.Timeout = 1 * time.Second

	compressionPlugin := plugin.NewCompressionPlugin(rpcx.CompressSnappy)
	client.PluginContainer.Add(compressionPlugin)

	args := &Args{7, 8}
	var reply Reply

	for i := 0; i < 100; i++ {

		go func() {
			err := client.Call(context.Background(), "Arith.Mul", args, &reply)
			if err != nil {
				log.Infof("error for Arith: %d*%d, %v", args.A, args.B, err)
			} else {
				log.Infof("Arith: %d*%d=%d", args.A, args.B, reply.C)
			}
		}()

		time.Sleep(2 * time.Second)
	}

	select {}

	client.Close()
}
