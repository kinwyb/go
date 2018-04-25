package rpcx

import "github.com/kinwyb/go/generate"

func init() {
	generate.RegisterLayouter("rpcx", &lay{})
}
