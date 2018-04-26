package rpcx

import "github.com/kinwyb/go/generate"

func init() {
	generate.RegisterLayouter("rpcxclient", &layclient{})
	generate.RegisterLayouter("rpcx", &lay{})
}
