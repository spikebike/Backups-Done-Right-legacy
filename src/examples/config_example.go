package main

import (
	"fmt"
	"github.com/kless/goconfig/config"
)

func main() {
	c, _ := config.ReadDefault("../../config.cfg")

	client_public_key, _ := c.String("Client", "public_key")
	server_public_key, _ := c.String("Server", "public_key")
	server_max_cores, _ := c.Int("Server", "max_cores")

	fmt.Println("client public key: ", client_public_key)
	fmt.Println("server public key: ", server_public_key)
	fmt.Println("server max cores: ", server_max_cores)
}
