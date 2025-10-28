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
type handlerLoggedIn func(*state, command, database.User) error

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

	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		log.Printf("Usage: %s <name> <url>\n", cmd.name)
		return fmt.Errorf("feed name and url is required")
	}

	name := cmd.args[0]
	if len(cmd.args) == 1 {
		log.Printf("Usage: %s %s <url>\n", cmd.name, name)
		return fmt.Errorf("url is required")
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

	follow, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println("=====================================")
	fmt.Println("Feed follow created successfully:")
	printCreateFeedFollow(follow)
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		log.Printf("Usage: %s <url>\n", cmd.name)
		return fmt.Errorf("url is required")
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	follow, err := s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}
	fmt.Println("Feed follow created successfully:")
	printCreateFeedFollow(follow)
	fmt.Println("=====================================")
	return nil
}

func handlerGetFollows(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't feed follows: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("No feed follows found.")
		return nil
	}

	fmt.Printf("Found %d feed follows:\n", len(follows))
	for _, feed := range follows {
		printGetFeedFollows(feed)
		fmt.Println("=====================================")
	}
	return nil
}

func printCreateFeedFollow(ffw database.CreateFeedFollowRow) {
	printFeedFollows(
		ffw.ID, ffw.FeedID, ffw.UserID,
		ffw.CreatedAt, ffw.UpdatedAt,
		ffw.FeedName, ffw.UserName,
	)
}

func printGetFeedFollows(ffwu database.GetFeedFollowsForUserRow) {
	printFeedFollows(
		ffwu.ID, ffwu.FeedID, ffwu.UserID,
		ffwu.CreatedAt, ffwu.UpdatedAt,
		ffwu.FeedName, ffwu.UserName,
	)
}

func printFeedFollows(id, fid, uid uuid.UUID, ca, ua time.Time, fname, uname string) {
	fmt.Printf("* ID:            %s\n", id)
	fmt.Printf("* Created:       %v\n", ca)
	fmt.Printf("* Updated:       %v\n", ua)
	fmt.Printf("* UserID:        %v\n", uid)
	fmt.Printf("* FeedID:        %v\n", fid)
	fmt.Printf("* FeedName:      %s\n", fname)
	fmt.Printf("* UserName:      %s\n", uname)
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteAllUsers(context.Background()); err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	fmt.Println("Database reset successfully!")
	return nil
}
