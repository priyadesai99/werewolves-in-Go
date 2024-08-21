package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"werewolves-go/types"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
)

type client struct {
	username  string
	serverPID *actor.PID
	logger    *slog.Logger
}

func newClient(username string, serverPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &client{
			username:  username,
			serverPID: serverPID,
			logger:    slog.Default(),
		}
	}
}

func (c *client) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case *types.Message:
		fmt.Printf("%s: %s\n", msg.Username, msg.Msg)
	case actor.Started:
		ctx.Send(c.serverPID, &types.Connect{
			Username: c.username,
		})
	case actor.Stopped:
		c.logger.Info("client stopped")
	}
}

// Handles keyboard interrupts from the client
// https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html
func getFireSignalsChannel() chan os.Signal {

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGHUP,  // "terminal is disconnected"
	)
	return c
}

// client cleanup service
func cleanup(serverPID *actor.PID, clientPID *actor.PID, e *actor.Engine) {
	e.SendWithSender(serverPID, &types.Disconnect{}, clientPID)
	e.Poison(clientPID).Wait()
	slog.Info("client disconnected")
}

func main() {
	var (
		listenAt  = flag.String("listen", "", "specify address to listen to, will pick a random port if not specified")
		connectTo = flag.String("connect", "127.0.0.1:4000", "the address of the server to connect to")
		username  = flag.String("username", "", "Enter username for client")
	)
	flag.Parse()

	for {
		if len(*username) == 0 {
			slog.Error("Username cannot be empty")
			slog.Info("Please enter username: ")
			fmt.Scan(username)
			continue
		}
		break
	}

	if *listenAt == "" {
		*listenAt = fmt.Sprintf("127.0.0.1:%d", rand.Int31n(50000)+10000)
	}
	rem := remote.New(*listenAt, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(rem))
	if err != nil {
		slog.Error("failed to create engine", "err", err)
		os.Exit(1)
	}

	var (
		// the process ID of the server
		serverPID = actor.NewPID(*connectTo, "server/primary")
		// Spawn our client receiver
		clientPID = e.Spawn(newClient(*username, serverPID), "client", actor.WithID(*username))
		scanner   = bufio.NewScanner(os.Stdin)
	)

	// Interrupt handling
	exitChan := getFireSignalsChannel()
	go func() {
		<-exitChan
		cleanup(serverPID, clientPID, e)
		os.Exit(1)
	}()

	fmt.Println("Type 'quit' and press return to exit.")
	for scanner.Scan() {
		msg := &types.Message{
			Msg:      scanner.Text(),
			Username: *username,
		}
		// We use SendWithSender here so the server knows who
		// is sending the message.
		if msg.Msg == "quit" {
			cleanup(serverPID, clientPID, e)
			break
		}
		e.SendWithSender(serverPID, msg, clientPID)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("failed to read message from stdin", "err", err)
	}
}
