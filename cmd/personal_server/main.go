package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"

	"github.com/nolanjannotta/personal_go_server/pkg/httpServer"
	"github.com/nolanjannotta/personal_go_server/pkg/sshApp"
)

func main() {

	http_server := httpServer.SetUp()
	ssh_app := sshApp.SetUp()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go sshApp.Start(ssh_app, done)

	go httpServer.Start(http_server, done)

	<-done

	log.Info("Stopping servers")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()

	if err := ssh_app.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
	if err := http_server.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}

}
