package graph

import (
	"context"
	"encoding/json"
	"fmt"

	_ "github.com/lib/pq"
	db "github.com/mattmazer1/graphql-api/database"
	"github.com/mattmazer1/graphql-api/graph/model"
	auth "github.com/mattmazer1/graphql-api/middleware"
	nat "github.com/mattmazer1/graphql-api/nats"
	"github.com/mattmazer1/graphql-api/utils"
	nats "github.com/nats-io/nats.go"
)

// CreatePlayer is the resolver for the createPlayer field.
func (r *mutationResolver) CreatePlayer(ctx context.Context, player model.InputPlayer) (*model.Player, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}

	createdPlayer := &model.Player{
		Pos:        player.Pos,
		Name:       player.Name,
		Age:        player.Age,
		Experience: player.Experience,
		Stats: &model.Stats{
			Season:    player.Stats.Season,
			Points:    player.Stats.Points,
			ThreePt:   player.Stats.ThreePt,
			Rebounds:  player.Stats.Rebounds,
			Assists:   player.Stats.Assists,
			Steals:    player.Stats.Steals,
			Blocks:    player.Stats.Blocks,
			TurnOvers: player.Stats.TurnOvers,
			Mp:        player.Stats.Mp},
	}

	dbErr := db.CreatePlayer(ctx, player)

	if dbErr != nil {
		return nil, fmt.Errorf("could not delete player: %w", dbErr)
	}

	createdPlayerJSON, err := json.Marshal(createdPlayer)

	if err != nil {
		return nil, fmt.Errorf("could not marshal created player: %w", err)
	}

	nat.Nc.Publish("create", createdPlayerJSON)

	return createdPlayer, nil
}

// UpdatePlayer is the resolver for the updatePlayer field.
func (r *mutationResolver) UpdatePlayer(ctx context.Context, player model.InputUpdatePlayer) (*model.Player, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}

	updatedPlayer := &model.Player{
		Pos:        *player.Pos,
		Name:       *player.Name,
		Age:        *player.Age,
		Experience: *player.Experience,
		Stats: &model.Stats{
			Season:    *player.Stats.Season,
			Points:    *player.Stats.Points,
			ThreePt:   *player.Stats.ThreePt,
			Rebounds:  *player.Stats.Rebounds,
			Assists:   *player.Stats.Assists,
			Steals:    *player.Stats.Steals,
			Blocks:    *player.Stats.Blocks,
			TurnOvers: *player.Stats.TurnOvers,
			Mp:        *player.Stats.Mp},
	}

	err := db.UpdatePlayer(ctx, player)

	if err != nil {
		return nil, fmt.Errorf("could not update player: %w", err)
	}

	return updatedPlayer, nil
}

// DeletePlayer is the resolver for the deletePlayer field.
func (r *mutationResolver) DeletePlayer(ctx context.Context, name string) (*model.Player, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("access denied")
	}

	err := db.DeletePlayer(ctx, name)

	if err != nil {
		return nil, fmt.Errorf("could not delete player: %w", err)
	}

	return &model.Player{Name: name}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.InputUser) (string, error) {
	user := &model.User{
		Username: input.Username,
		Password: input.Password,
	}

	correct, err := db.Authenticate(user)

	if err != nil {
		return "", fmt.Errorf("could not authenticate user %w", err)
	}
	if !correct {
		return "", fmt.Errorf("could not authenticate user")
	}

	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil

}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, token string) (string, error) {
	username, err := utils.ParseToken(token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}

	newToken, err := utils.GenerateToken(username)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}
	return newToken, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.InputUser) (string, error) {
	err := db.CreateUsr(ctx, input)

	if err != nil {
		return "", fmt.Errorf("could no create user to db: %w", err)
	}

	token, err := utils.GenerateToken(input.Username)

	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUsername(ctx context.Context, user model.UpdateUsername) (string, error) {
	userauth := auth.ForContext(ctx)
	if userauth == nil {
		return "", fmt.Errorf("access denied")
	}

	id, err := db.GetUserId(user.OldUsername)

	if err != nil {
		return "", fmt.Errorf("could no get user id: %w", err)
	}

	user.ID = &id

	dbErr := db.UpdateUsername(ctx, user)

	if dbErr != nil {
		return "", fmt.Errorf("could not update username: %w", err)
	}

	return fmt.Sprintf("Successfully updated username for user id - %s", *user.ID), nil
}

func (r *mutationResolver) UpdatePassword(ctx context.Context, user model.UpdatePassword) (string, error) {
	userauth := auth.ForContext(ctx)
	if userauth == nil {
		return "", fmt.Errorf("access denied")
	}

	id, err := db.GetUserId(user.Username)

	if err != nil {
		return "", fmt.Errorf("could no get user id: %w", err)
	}

	user.ID = &id

	dbErr := db.UpdatePassword(ctx, user)

	if dbErr != nil {
		return "", fmt.Errorf("could not update password: %w", err)
	}

	return fmt.Sprintf("Successfully updated password for user id - %s", *user.ID), nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, username string) (string, error) {
	userauth := auth.ForContext(ctx)
	if userauth == nil {
		return "", fmt.Errorf("access denied")
	}

	err := db.DeleteUser(ctx, username)

	if err != nil {
		return "", fmt.Errorf("could not delete user: %w", err)
	}

	return fmt.Sprintf("Deleted user %s", username), nil
}

// Player is the resolver for the player field.
func (r *queryResolver) Player(ctx context.Context, name string) (*model.Player, error) {
	player, err := db.GetPlayer(ctx, name)

	if err != nil {
		return nil, fmt.Errorf("could not get player: %w", err)
	}

	return player, nil
}

// GetUserID is the resolver for the getUserId field.
func (r *queryResolver) GetUserID(ctx context.Context, username string) (string, error) {
	id, err := db.GetUserId(username)

	if err != nil {
		return "", fmt.Errorf("could not get user id: %w", err)
	}

	return id, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, username string) (*model.User, error) {
	user, err := db.GetUser(username)

	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	return user, nil
}

// Player is the resolver for the player field.
func (r *subscriptionResolver) Player(ctx context.Context) (<-chan *model.Player, error) {
	ch := make(chan *model.Player)

	go func() {

		sub, err := nat.Nc.Subscribe("create", func(m *nats.Msg) {
			player := &model.Player{}

			err := json.Unmarshal(m.Data, player)
			if err != nil {
				fmt.Println("could not umarshall created player: %w", err)
				return
			}

			ch <- player
		})

		if err != nil {
			close(ch)
			fmt.Println("could not subscribe to created player: %w", err)
			return
		}

		go func() {
			<-ctx.Done()
			sub.Unsubscribe()
		}()

	}()

	return ch, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
