// Package gogator contains the implementaion of gogator
// cli app.
package gogator

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/prchop/gogator/internal/config"
	"github.com/prchop/gogator/internal/database"
)

type state struct {
	cfg *config.Settings
	db  *database.Queries
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

func Run() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v\n", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}
	defer db.Close()
	dbq := database.New(db)

	programState := &state{cfg: &cfg, db: dbq}
	cmds := &commands{registeredCmds: make(map[string]handler)}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("users", handlerListUsers)
	cmds.register("reset", handlerReset)
	cmds.register("agg", handlerAggregate)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerListFeeds)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{
		name: cmdName, args: cmdArgs,
	})
	if err != nil {
		log.Fatalf("[ERROR]: %v", err)
	}
}
