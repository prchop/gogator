package gogator

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		log.Printf("Usage: %s <username>\n", cmd.name)
		return fmt.Errorf("username is required")
	}

	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user %q: %w", name, err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User %q switched successfully!\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		log.Printf("Usage: %s <username>\n", cmd.name)
		return fmt.Errorf("username is required")
	}

	name := cmd.args[0]
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      name,
		})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User %q created successfully:\n", name)
	printUser(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}

func printUser(u database.User) {
	fmt.Printf(" * ID:      %v\n", u.ID)
	fmt.Printf(" * Name:    %v\n", u.Name)
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
	cmds.register("reset", handlerReset)

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
