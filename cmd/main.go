package main

import (
	"github.com/fabienbellanger/fiber-boilerplate/pkg/infrastructure/cli"
	"log"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatalln(err)
	}
}
