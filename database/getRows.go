package db

import (
	"database/sql"
	"fmt"

	"github.com/mattmazer1/graphql-api/graph/model"
)

func getRows(rows *sql.Rows, name string) (*model.Player, error) {
	var position string
	var age int
	var experience int
	var season string
	var points float64
	var threept float64
	var rebounds float64
	var assists float64
	var steals float64
	var blocks float64
	var turnovers float64
	var mp float64

	var player *model.Player

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&name,
			&position,
			&age,
			&experience,
			&season,
			&points,
			&threept,
			&rebounds,
			&assists,
			&steals,
			&blocks,
			&turnovers,
			&mp,
		); err != nil {
			return nil, fmt.Errorf("could not scan player stats: %w", err)
		}
		player = &model.Player{
			Pos:        model.Position(position),
			Name:       name,
			Age:        age,
			Experience: experience,
			Stats: &model.Stats{
				Season:    season,
				Points:    points,
				ThreePt:   threept,
				Rebounds:  rebounds,
				Assists:   assists,
				Steals:    steals,
				Blocks:    blocks,
				TurnOvers: turnovers,
				Mp:        mp},
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return player, nil
}

func getUserIdRows(rows *sql.Rows) (string, error) {
	var userId string
	var id string

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&id,
		); err != nil {
			return "", fmt.Errorf("could not scan user id: %w", err)
		}
		userId = id
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("rows error: %w", err)
	}

	return userId, nil
}

func getUserRows(rows *sql.Rows) (*model.User, error) {
	var id string
	var username string
	var password string

	var user *model.User

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&id, &username, &password,
		); err != nil {
			return nil, fmt.Errorf("could not scan user: %w", err)
		}
		user = &model.User{
			ID:       id,
			Password: password,
			Username: username,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return user, nil
}

func getUserPassword(rows *sql.Rows) (string, error) {
	var password string
	var hashedPassword string

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&password,
		); err != nil {
			return "", fmt.Errorf("could not scan user info: %w", err)
		}
		hashedPassword = password
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("rows error: %w", err)
	}

	return hashedPassword, nil
}
