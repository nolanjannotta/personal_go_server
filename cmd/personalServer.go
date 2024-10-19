package personalServer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/joho/godotenv"

	"github.com/nolanjannotta/personal_go_server/pkg/httpServer"
	"github.com/nolanjannotta/personal_go_server/pkg/sshApp"
)

func Launch() {
	// godotenv.Load("../../.env")

	sess := session.Must(session.NewSession())
	fmt.Println(sess)

	godotenv.Load("./.env") // dockerfile

	// l, err := net.Listen("tcp", ":42069")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// m := cmux.New(l)

	http_server := httpServer.SetUp()
	ssh_app := sshApp.SetUp()

	// httpL := m.Match(cmux.HTTP1Fast())
	// trpcL := m.Match(cmux.Any())

	// go http_server.Serve(httpL)
	// go ssh_app.Serve(trpcL)

	// log.Info("Starting servers")
	// m.Serve()

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
