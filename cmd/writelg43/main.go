package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	oplg43ud79 "github.com/sisisin/op-lg43ud79"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := oplg43ud79.RunWriteLG43(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
