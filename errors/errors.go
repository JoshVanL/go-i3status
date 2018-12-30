package errors

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func NewSignalHandler() <-chan int {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGCONT, syscall.SIGSTOP, syscall.SIGKILL, syscall.SIGINT)

	ch := make(chan int)

	go func() {
		for s := range sig {
			switch s {
			case syscall.SIGCONT:
				ch <- 1
				break

			case syscall.SIGSTOP:
				ch <- 0
				break

			case syscall.SIGKILL, syscall.SIGINT:
				ch <- -1

			default:
				break
			}
		}
	}()

	return ch
}

func Kill(err error) {
	f, ferr := os.OpenFile("/var/run/go-i3status/err.log", os.O_CREATE|os.O_WRONLY, 0666)
	if ferr == nil {
		fmt.Fprint(f, err)
		f.Close()
	}

	fmt.Fprint(os.Stderr, err)

	os.Exit(1)
}
