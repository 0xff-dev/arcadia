package client

import (
	"sync"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var (
	once sync.Once
	c    dynamic.Interface
)

func InitClient(cfg *rest.Config) {
	once.Do(func() {
		c = dynamic.NewForConfigOrDie(cfg)
	})
}

func GetClient() dynamic.Interface {
	return c
}
