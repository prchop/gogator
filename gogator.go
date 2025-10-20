package gogator

import (
	"fmt"
	"log"
	"os"

	"github.com/prchop/gogator/internal/config"
)

type state struct {
	config *config.Settings
}

type command struct {
	name string
	args []string
}

type handler func(*state, command) error

type commands struct {
	registeredCmds map[string]handler
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCmds[cmd.name]
	if !ok {
		return fmt.Errorf("command not found")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f handler) {
	c.registeredCmds[name] = f
}

func loginHandler(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		log.Printf("Usage: %s <username>\n", cmd.name)
		return fmt.Errorf("username is required")
	}

	name := cmd.args[0]
	if err := s.config.SetUser(name); err != nil {
		return fmt.Errorf("couldn't set current user: %v", err)
	}
	fmt.Printf("user has been set to %q\n", name)
	return nil
}

func Run() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}

	programState := &state{config: &cfg}
	cmds := &commands{registeredCmds: make(map[string]handler)}
	cmds.register("login", loginHandler)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{
		name: cmdName, args: cmdArgs,
	})
	if err != nil {
		log.Fatal(err)
	}
}
