package db

import (
	"context"
	"fmt"

	"github.com/mattmazer1/graphql-api/graph/model"
	"github.com/mattmazer1/graphql-api/utils"
)

func GetPlayer(ctx context.Context, name string) (*model.Player, error) {
	rows, err := Db.QueryContext(ctx, `SELECT * FROM players
	WHERE name = $1;`, name)

	if err != nil {
		return nil, fmt.Errorf("could not get player: %w", err)
	}

	player, err := getRows(rows, name)

	if err != nil {
		return nil, fmt.Errorf("could not not get player rows: %w", err)
	}

	return player, nil
}

func DeletePlayer(ctx context.Context, name string) error {
	_, err := Db.ExecContext(ctx, `DELETE FROM players WHERE name = $1;`,
		name,
	)

	if err != nil {
		return fmt.Errorf("could not delete player: %w", err)
	}
	return nil
}

func CreatePlayer(ctx context.Context, player model.InputPlayer) error {
	_, err := Db.ExecContext(ctx, `INSERT INTO players (
			name,
			position,
			age,
			experience,
			season,
			points,
			threept,
			rebounds,
			assists,
			steals,
			blocks,
			turnovers,
			mp)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		player.Name,
		player.Pos,
		player.Age,
		player.Experience,
		player.Stats.Season,
		player.Stats.Points,
		player.Stats.ThreePt,
		player.Stats.Rebounds,
		player.Stats.Assists,
		player.Stats.Steals,
		player.Stats.Blocks,
		player.Stats.TurnOvers,
		player.Stats.Mp,
	)

	if err != nil {
		return fmt.Errorf("could not create player stats: %w", err)
	}

	return nil
}

func UpdatePlayer(ctx context.Context, player model.InputUpdatePlayer) error {
	_, err := Db.ExecContext(ctx,
		`UPDATE players
		SET
		name = $1,
		position = $2,
		age = $3,
		experience = $4,
		season = $5,
		points = $6,
		threept = $7,
		rebounds = $8,
		assists = $9,
		steals = $10,
		blocks = $11,
		turnovers = $12,
		mp = $13
		WHERE name = $1
		`,
		player.Name,
		player.Pos,
		player.Age,
		player.Experience,
		player.Stats.Season,
		player.Stats.Points,
		player.Stats.ThreePt,
		player.Stats.Rebounds,
		player.Stats.Assists,
		player.Stats.Steals,
		player.Stats.Blocks,
		player.Stats.TurnOvers,
		player.Stats.Mp,
	)

	if err != nil {
		return fmt.Errorf("could not update player stats: %w", err)
	}

	return nil
}

func GetUserId(username string) (string, error) {
	rows, err := Db.Query(`SELECT id FROM users
	WHERE username = $1;`, username)

	if err != nil {
		return "", fmt.Errorf("could not get user: %w", err)
	}

	id, err := getUserIdRows(rows)

	if err != nil {
		return "", fmt.Errorf("could not not get user rows: %w", err)
	}

	return id, nil
}

func GetUser(username string) (*model.User, error) {
	rows, err := Db.Query(`SELECT * FROM users
	WHERE username = $1;`, username)

	if err != nil {
		return nil, fmt.Errorf("could not get user: %w", err)
	}

	user, err := getUserRows(rows)

	if err != nil {
		return nil, fmt.Errorf("could not not get user rows: %w", err)
	}

	return user, nil
}

func CreateUsr(ctx context.Context, user model.InputUser) error {
	hashedPassword, pswerr := utils.HashPassword(user.Password)
	if pswerr != nil {
		return fmt.Errorf("could not hash user password: %w", pswerr)
	}
	_, err := Db.ExecContext(ctx, `INSERT INTO users (
		id,
		username,
		password
	)
	VALUES (gen_random_uuid(), $1, $2)`,
		user.Username,
		hashedPassword,
	)

	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}

	return nil
}

func UpdateUsername(ctx context.Context, user model.UpdateUsername) error {
	_, err := Db.ExecContext(ctx, `
	UPDATE users
		SET
		username = $2
		WHERE
		id = $1`,
		user.ID,
		user.NewUsername,
	)

	if err != nil {
		return fmt.Errorf("could not update username: %w", err)
	}

	return nil
}

func UpdatePassword(ctx context.Context, user model.UpdatePassword) error {
	hashedPassword, pswderr := utils.HashPassword(user.NewPassword)
	if pswderr != nil {
		return fmt.Errorf("could not hash user password: %w", pswderr)
	}

	_, err := Db.ExecContext(ctx, `
	UPDATE users
		SET
		password = $2
		WHERE
		id = $1`,
		user.ID,
		hashedPassword,
	)

	if err != nil {
		return fmt.Errorf("could not update user password: %w", err)
	}

	return nil
}

func DeleteUser(ctx context.Context, username string) error {
	_, err := Db.ExecContext(ctx, `DELETE FROM users WHERE username = $1;`,
		username,
	)

	if err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}

	return nil
}

func Authenticate(user *model.User) (bool, error) {
	rows, err := Db.Query(`SELECT password FROM users
	WHERE username = $1;`, user.Username)

	if err != nil {
		return false, fmt.Errorf("could not get user: %w", err)
	}

	dbHashedPassword, err := getUserPassword(rows)

	if err != nil {
		return false, fmt.Errorf("could not not get player rows: %w", err)
	}

	return utils.CheckPasswordHash(user.Password, dbHashedPassword), nil
}
