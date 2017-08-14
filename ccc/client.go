package main

import (
	"context"
	"time"

	"github.com/smallnest/pool"
	"github.com/smallnest/rpcx"
	"github.com/smallnest/rpcx/log"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

func main() {
	s := &rpcx.DirectClientSelector{Network: "quic", Address: "127.0.0.1:8972", DialTimeout: 10 * time.Second}

	clientPool := &pool.Pool{
		New: func() interface{} {
			return rpcx.NewClient(s)
		},
	}

	client := clientPool.Get().(*rpcx.Client)
	client.Timeout = 3 * time.Second

	args := &Args{7, 8}
	var reply Reply
	err := client.Call(context.Background(), "Arith.Mul", args, &reply)
	if err != nil {
		log.Infof("error for Arith: %d*%d, %v", args.A, args.B, err)
	} else {
		log.Infof("Arith: %d*%d=%d, client: %p", args.A, args.B, reply.C, client)
	}

	rpcx.Reconnect = nil

	clientPool.Put(client)

	//restart server
	time.Sleep(20 * time.Second)

	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	// }()

	log.Infof("starting test")
	for i := 0; i < 10; i++ {
		client = clientPool.Get().(*rpcx.Client)
		client.Timeout = 3 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = client.Call(ctx, "Arith.Mul", args, &reply)
		if err != nil {
			log.Infof("error for Arith: %d*%d, %v", args.A, args.B, err)
		} else {
			log.Infof("Arith: %d*%d=%d, client: %p", args.A, args.B, reply.C, client)
		}

		clientPool.Put(client)
	}

	// client = clientPool.Get().(*rpcx.Client)
	// client.Timeout = 3 * time.Second

	// ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// err = client.Call(context.Background(), "Arith.Mul", args, &reply)
	// if err != nil {
	// 	log.Infof("error for Arith: %d*%d, %v", args.A, args.B, err)
	// } else {
	// 	log.Infof("Arith: %d*%d=%d, client: %p", args.A, args.B, reply.C, client)
	// }

	// clientPool.Put(client)

	clientPool.Range(func(v interface{}) bool {
		c := v.(*rpcx.Client)
		c.Close()
		return true
	})
	clientPool.Reset()
}
