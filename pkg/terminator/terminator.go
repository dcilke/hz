package terminator

import (
	"os"
	"os/signal"
	"syscall"
)

func OnSig(f func() int) {
	c := make(chan os.Signal, 3)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		os.Exit(f())
	}()
}
