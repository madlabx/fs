package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/madlabx/pkgx/log"
)

func enablePProfile(debugPort int) {
	if err := http.ListenAndServe(fmt.Sprintf(":%d", debugPort), nil); err != nil {
		log.Errorf("Cannot enable debug port")
	}
}

// signal.Notify(sigCh, syscall.SIGUSR1, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
func cleanupHandler(c chan os.Signal, cancel context.CancelFunc) { //nolint:interfacer
	defer func() {
		cancel()
		// TODO: how to know all cleanup has finished
		time.Sleep(3 * time.Second)
		os.Exit(0)
	}()
	for {
		sig := <-c
		switch sig {

		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL:
			log.Errorf("receive signal %v, program will exit", sig.String())
			return
		case syscall.SIGSEGV:
			log.Errorf("recive signal %v, callstack: %v", sig.String(), debug.Stack())
			return
		default:
			//log.Errorf("receive signal %v, ignore it", sig.String())
		}
	}
}
