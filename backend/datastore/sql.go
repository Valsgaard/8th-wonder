package datastore

import (
	"fmt"

	"github.com/jackc/pgx"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

type sqlDatabase struct {
	connection         *pgx.ConnPool
	preparedStatements map[string]*pgx.PreparedStatement
}

var _ port.Datastore = &sqlDatabase{}

func (db *sqlDatabase) Prepare(name, query string) common.Error {
	_, ok := db.preparedStatements[name]
	if ok {
		return nil
	}

	ps, err := db.connection.Prepare(name, query)
	if err != nil {
		return common.NewError(err, "")
	}

	db.preparedStatements[name] = ps
	return nil
}

// NewSQLConnection creates a connection to a PostgreSQL server
func NewSQLConnection(uri string, maxConnections int) (port.Datastore, common.Error) {
	config, err := pgx.ParseURI(uri)
	if err != nil {
		fmt.Println("Hi!")
		return nil, common.NewError(err, "")
	}

	if maxConnections < 2 {
		maxConnections = 5 // pgx default
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: maxConnections,
	}

	conn, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		fmt.Println("Hi2!")
		return nil, common.NewError(err, "")
	}

	db := &sqlDatabase{
		connection:         conn,
		preparedStatements: make(map[string]*pgx.PreparedStatement),
	}

	return db, nil
}

// NewUser ...
func (db *sqlDatabase) NewUser(id, name string) (*port.User, common.Error) {
	qName := "newUser"
	q := `INSERT INTO users (id, name) VALUES($1, $2) RETURNING id, name;`

	if err := db.Prepare(qName, q); err != nil {
		return nil, err
	}

	user := new(port.User)
	err := db.connection.QueryRow(qName, id, name).
		Scan(&user.UserID, &user.Name)

	if err != nil {
		return nil, common.NewError(err, "")
	}

	return user, nil
}

// GetUsers ...
func (db *sqlDatabase) GetUsers() ([]*port.User, common.Error) {
	qName := "getUser"
	q := `SELECT id, name FROM users ORDER BY id;`

	if err := db.Prepare(qName, q); err != nil {
		return nil, err
	}

	rows, err := db.connection.Query(qName)
	if err != nil {
		return nil, common.NewError(err, "")
	}

	users := make([]*port.User, 0)
	for rows.Next() {
		user := new(port.User)

		if err := rows.Scan(&user.UserID, &user.Name); err != nil {
			return nil, common.NewError(err, "")
		}

		users = append(users, user)
	}

	return users, nil
}

func (db *sqlDatabase) UserExists(id string) (bool, common.Error) {
	qName := "userExists"
	q := `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1);`

	if err := db.Prepare(qName, q); err != nil {
		return false, err
	}

	var check bool
	err := db.connection.QueryRow(qName, id).
		Scan(&check)

	if err != nil {
		return false, common.NewError(err, "")
	}

	return check, nil
}

// UpdateGameState ...
func (db *sqlDatabase) UpdateGameState(userID string, gamesPlayed, score int) common.Error {
	qName := "updateGameState"
	q := `UPDATE users SET (games_played, score) = ($1, $2) WHERE id = $3;`

	if err := db.Prepare(qName, q); err != nil {
		return err
	}

	_, err := db.connection.Exec(qName, gamesPlayed, score, userID)
	if err != nil {
		return common.NewError(err, "")
	}

	return nil
}

// GetGameState ...
func (db *sqlDatabase) GetGameState(userID string) (*port.GameState, common.Error) {
	qName := "getGameState"
	q := `SELECT (games_played, score) FROM users WHERE id = $1;`

	if err := db.Prepare(qName, q); err != nil {
		return nil, err
	}

	gameState := new(port.GameState)
	err := db.connection.QueryRow(qName, userID).Scan(&gameState.GamesPlayed, &gameState.Score)
	if err != nil {
		return nil, common.NewError(err, "")
	}

	return gameState, nil
}

// UpdateFriends ...
func (db *sqlDatabase) UpdateFriends(userID string, friends []string) common.Error {
	qName := "updateFriends"
	q := `UPDATE users SET (friends) = ($1) WHERE id = $2;`

	if err := db.Prepare(qName, q); err != nil {
		return err
	}

	_, err := db.connection.Exec(qName, friends, userID)
	if err != nil {
		return common.NewError(err, "")
	}

	return nil
}

// GetFriends ...
func (db *sqlDatabase) GetFriends(userID string) ([]*port.Friend, common.Error) {
	qName := "getFriends"
	q := `SELECT id, name, score FROM users WHERE id = ANY((SELECT unnest(friends) FROM users WHERE id = $1)) ORDER BY id;`

	if err := db.Prepare(qName, q); err != nil {
		return nil, err
	}

	rows, err := db.connection.Query(qName, userID)
	if err != nil {
		return nil, common.NewError(err, "")
	}

	friends := make([]*port.Friend, 0)
	for rows.Next() {
		friend := new(port.Friend)

		if err := rows.Scan(&friend.UserID, &friend.Name, &friend.HighScore); err != nil {
			return nil, common.NewError(err, "")
		}

		friends = append(friends, friend)
	}

	return friends, nil
}

// DeleteUser ...
func (db *sqlDatabase) DeleteUser(userID string) common.Error {
	qName := "deleteUser"
	q := `DELETE users WHERE id = $1;`

	if err := db.Prepare(qName, q); err != nil {
		return err
	}

	_, err := db.connection.Exec(qName, userID)
	if err != nil {
		return common.NewError(err, "")
	}

	return nil
}
