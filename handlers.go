package gogator

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/prchop/gogator/internal/database"
)

type handler func(*state, command) error

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

func handlerListUsers(s *state, cmd command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get all users: %w", err)
	}
	for _, u := range users {
		if u.Name == s.cfg.UserName {
			fmt.Printf(" * Name:    %v (current)\n", u.Name)
			continue
		}
		fmt.Printf(" * Name:    %v\n", u.Name)
	}
	return nil
}

func printUser(u database.User) {
	fmt.Printf(" * ID:      %v\n", u.ID)
	fmt.Printf(" * Name:    %v\n", u.Name)
}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		log.Printf("Usage: %s <url>\n", cmd.name)
		return fmt.Errorf("url is required")
	}

	feed, err := fetchFeed(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	fmt.Printf("Feed: %+v", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		log.Printf("Usage: %s <name> <url>\n", cmd.name)
		return fmt.Errorf("feed name and url is required")
	}

	name := cmd.args[0]
	if len(cmd.args) == 1 {
		log.Printf("Usage: %s %s <url>\n", cmd.name, name)
		return fmt.Errorf("url is required")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.UserName)
	if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	}

	url := cmd.args[1]
	feed, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      name,
			Url:       url,
			UserID:    user.ID,
		})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println("=====================================")
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("=====================================")
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}
