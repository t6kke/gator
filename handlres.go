package main

import (
	"fmt"
	"time"
	"context"
	"github.com/google/uuid"
	"github.com/t6kke/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No username provided, username is required --- Usage: %s <name>", cmd.name)
	}

	new_ctx := context.Background()
	_, err := s.dbq.GetUser(new_ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = s.conf.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User '%s' has been successfully configured for session\n", cmd.args[0])
	return nil
}



func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No username provided, username is required --- Usage: %s <name>", cmd.name)
	}

	new_uuid := uuid.New()
	current_time := time.Now()
	user_name := cmd.args[0]

	new_user := database.CreateUserParams{
		ID:        new_uuid,
		CreatedAt: current_time,
		UpdatedAt: current_time,
		Name:      user_name,
	}

	new_ctx := context.Background()
	user, err := s.dbq.GetUser(new_ctx, user_name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return fmt.Errorf("%w", err)
	}

	if user.Name != "" {
		fmt.Errorf("User with name already exists")
	}

	user, err = s.dbq.CreateUser(new_ctx, new_user)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = s.conf.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User '%s' has been successfully added to the database\n", user_name)
	fmt.Printf("DEBUG --- uuid: '%v' --- timestamp: '%v' --- user: '%s'\n", new_user.ID, new_user.CreatedAt, new_user.Name)

	return nil
}



func handlerReset(s *state, cmd command) error {
	new_ctx := context.Background()
	err := s.dbq.DeleteAllUsers(new_ctx)
	return err
}



func handlerUsers(s *state, cmd command) error {
	new_ctx := context.Background()
	users, err := s.dbq.GetUsers(new_ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, user := range users {
		if user.Name == s.conf.Current_user_name {
			fmt.Println(user.Name, "(current)")
		} else {
			fmt.Println(user.Name)
		}
	}
	return nil
}


//just initial setup to confim that retreiving content is working as expected
func handlerAgg(s *state, cmd command) error {
	test_url := "https://www.wagslane.dev/index.xml"

	new_ctx := context.Background()
	feed, err := fetchFeed(new_ctx, test_url)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println(feed.Channel.Link)
	fmt.Println(feed.Channel.Title)
	fmt.Println(feed.Channel.Description)
	fmt.Println("-------------------------------------------------------")
	fmt.Println("items:")
	for i, item := range feed.Channel.Item {
		fmt.Printf("Item: %d --- Link: %s\n", i+1, item.Link)
		fmt.Println(item.Title)
		fmt.Println(item.Description)
		fmt.Println("-------------------------------------------------------")
	}

	return nil
}


func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("Missing argumets --- Usage: %s <name> <url>", cmd.name)
	}

	new_ctx := context.Background()

	user_uuid := user.ID
	feed_uuid := uuid.New()
	current_time := time.Now()
	feed_name := cmd.args[0]
	feed_url := cmd.args[1]

	new_feed := database.CreateFeedParams{
		ID:        feed_uuid,
		CreatedAt: current_time,
		UpdatedAt: current_time,
		Name:      feed_name,
		Url:       feed_url,
		UserID:    user_uuid,
	}

	feed, err := s.dbq.CreateFeed(new_ctx, new_feed)
	if err != nil {
		return err
	}

	fmt.Printf("Feed '%s' with url: '%s' has been successfully added to the database\n", feed_name, feed_url)
	fmt.Println("DEBUG --- ", feed)

	//including also following feed logic here, maybe this can be separated out into it's own function
	new_feed_follow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: current_time,
		UpdatedAt: current_time,
		UserID:    user_uuid,
		FeedID:    feed_uuid,
	}
	feed_follow, err := s.dbq.CreateFeedFollow(new_ctx, new_feed_follow)
	if err != nil {
		return err
	}
	fmt.Printf("Feed '%s' is successfully followed by: '%s'\n", feed_follow.FeedName, feed_follow.UserName)

	return nil
}


func handlerFeeds(s *state, cmd command) error {
	new_ctx := context.Background()
	feeds, err := s.dbq.GetFeeds(new_ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf("Name: %s\nURL: %s\ncreated by: %s\n--------------\n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No url provided, url is required --- Usage: %s <url>", cmd.name)
	}

	new_ctx := context.Background()
	feed_url := cmd.args[0]
	feed, err := s.dbq.GetFeed(new_ctx, feed_url)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return fmt.Errorf("%w", err)
	}
	if err != nil {
		return fmt.Errorf("Feed with url '%s' not found in database", feed_url)
	}

	new_uuid := uuid.New()
	current_time := time.Now()
	user_uuid := user.ID
	feed_uuid := feed.ID

	new_feed_follow := database.CreateFeedFollowParams{
		ID:        new_uuid,
		CreatedAt: current_time,
		UpdatedAt: current_time,
		UserID:    user_uuid,
		FeedID:    feed_uuid,
	}

	feed_follow, err := s.dbq.CreateFeedFollow(new_ctx, new_feed_follow)
	if err != nil {
		return err
	}

	fmt.Printf("Feed '%s' is successfully followed by: '%s'\n", feed_follow.FeedName, feed_follow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	new_ctx := context.Background()
	current_user := s.conf.Current_user_name

	user_uuid := user.ID
	follows, err := s.dbq.GetFeedFollowsForUser(new_ctx, user_uuid)
	if err != nil {
		return err
	}

	fmt.Printf("User '%s' is following:\n", current_user)
	for _, name := range follows {
		fmt.Println(name)
	}

	return nil
}




func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		new_ctx := context.Background()
		current_user := s.conf.Current_user_name
		user, err := s.dbq.GetUser(new_ctx, current_user)
		if err != nil && err.Error() != "sql: no rows in result set" {
			return fmt.Errorf("%w", err)
		}
		if err != nil {
			return fmt.Errorf("Current user '%s' not found in database", current_user)
		}
		return handler(s, cmd, user)
	}
}
