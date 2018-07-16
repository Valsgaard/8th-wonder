package port

import (
	"github.com/valsgaard/interview-case/backend/common"
)

// Datastore is the interface describing the methods required of an adepter
type Datastore interface {
	NewUser(id, name string) (*User, common.Error)
	GetUsers() ([]*User, common.Error)
	UserExists(id string) (bool, common.Error)

	UpdateGameState(userID string, gamesPlayed, score int) common.Error
	GetGameState(userID string) (*GameState, common.Error)

	UpdateFriends(userID string, friends []string) common.Error
	GetFriends(userID string) ([]*Friend, common.Error)

	// Used only for testing
	DeleteUser(userID string) common.Error
}

/**************************************************************************
***************************************************************************
**                                                                       **
**   Value Objects                                                       **
**                                                                       **
***************************************************************************
**************************************************************************/

// User is the value object used to input / output user related data from the adapter
type User struct {
	UserID string
	Name   string
}

// GameState is the value object used to input / output GameState related data from the adapter
type GameState struct {
	GamesPlayed int
	Score       int
}

// Friend is the value object used to input / output Friend related data from the adapter
type Friend struct {
	UserID    string
	Name      string
	HighScore int
}

/**************************************************************************
***************************************************************************
**                                                                       **
**   Error codes                                                         **
**                                                                       **
***************************************************************************
**************************************************************************/

// ErrEntryExists is used when an index key already exists
var ErrEntryExists = common.PrepareError("D001", "Entry already exists")

// ErrInvalidKey indicates that a given key is invalid or malformed
var ErrInvalidKey = common.PrepareError("D002", "Invalid key")
