package main

import (
	"log"

	"github.com/fabienbellanger/fiber-boilerplate/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatalln(err)
	}
}

