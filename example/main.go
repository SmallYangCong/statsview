package main

import (
	"os"
	"os/signal"

	"github.com/SmallYangCong/statsview"
)

func main() {
	mgr := statsview.New()

	go mgr.Start()

	sg := make(chan os.Signal, 1)
	signal.Notify(sg, os.Interrupt, os.Kill)
	<-sg
}
