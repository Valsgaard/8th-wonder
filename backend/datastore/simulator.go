package datastore

import (
	"sort"
	"sync"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

type datastoreSim struct {
	sync.Mutex

	Users map[string]*datastoreUser
}

var _ port.Datastore = &datastoreSim{}

// NewDatastoreSimulator creates an in-memory datastore, used for testing
func NewDatastoreSimulator() port.Datastore {
	return &datastoreSim{Users: make(map[string]*datastoreUser)}
}

type datastoreUser struct {
	userID    string
	name      string
	gameState port.GameState
	friendIDs []string
}

// NewUser ...
func (db *datastoreSim) NewUser(id, name string) (*port.User, common.Error) {
	db.Lock()
	defer db.Unlock()

	if _, ok := db.Users[id]; ok {
		return nil, common.NewError(port.ErrEntryExists, "User already exists")
	}

	// Create new user
	newUser := &datastoreUser{
		userID: id,
		name:   name,
		gameState: port.GameState{
			GamesPlayed: 0,
			Score:       0,
		},
		friendIDs: []string{},
	}

	db.Users[id] = newUser

	return &port.User{
		UserID: newUser.userID,
		Name:   newUser.name,
	}, nil
}

type userByUserID []*port.User

func (a userByUserID) Len() int           { return len(a) }
func (a userByUserID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a userByUserID) Less(i, j int) bool { return a[i].UserID < a[j].UserID }

// GetUsers ...
func (db *datastoreSim) GetUsers() ([]*port.User, common.Error) {
	users := make([]*port.User, 0, len(db.Users))
	for _, user := range db.Users {
		users = append(users, &port.User{UserID: user.userID, Name: user.name})
	}

	// Due to the randomization of maps, we'll have to sort it
	sort.Sort(userByUserID(users))

	return users, nil
}

func (db *datastoreSim) UserExists(id string) (bool, common.Error) {
	_, ok := db.Users[id]
	return ok, nil
}

// UpdateGameState ...
func (db *datastoreSim) UpdateGameState(userID string, gamesPlayed, score int) common.Error {
	db.Lock()
	defer db.Unlock()

	user, ok := db.Users[userID]
	if !ok {
		return common.NewError(port.ErrInvalidKey, "Invalid UserID")
	}

	user.gameState.GamesPlayed = gamesPlayed
	user.gameState.Score = score
	return nil
}

// GetGameState ...
func (db *datastoreSim) GetGameState(userID string) (*port.GameState, common.Error) {
	user, ok := db.Users[userID]
	if !ok {
		return nil, common.NewError(port.ErrInvalidKey, "Invalid UserID")
	}

	return &port.GameState{
		GamesPlayed: user.gameState.GamesPlayed,
		Score:       user.gameState.Score,
	}, nil
}

// UpdateFriends ...
func (db *datastoreSim) UpdateFriends(userID string, friends []string) common.Error {
	db.Lock()
	defer db.Unlock()

	user, ok := db.Users[userID]
	if !ok {
		return common.NewError(port.ErrInvalidKey, "Invalid UserID")
	}

	user.friendIDs = friends
	return nil
}

type friendByUserID []*port.Friend

func (a friendByUserID) Len() int           { return len(a) }
func (a friendByUserID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a friendByUserID) Less(i, j int) bool { return a[i].UserID < a[j].UserID }

// GetFriends ...
func (db *datastoreSim) GetFriends(userID string) ([]*port.Friend, common.Error) {
	user, ok := db.Users[userID]
	if !ok {
		return nil, common.NewError(port.ErrInvalidKey, "Invalid UserID")
	}

	friends := make([]*port.Friend, 0, len(user.friendIDs))
	for _, friendID := range user.friendIDs {
		friend, ok := db.Users[friendID]
		if !ok {
			// return nil, errors.New("Invalid userID")
			return nil, common.NewError(port.ErrInvalidKey, "Invalid UserID in friends")
		}

		friends = append(friends, &port.Friend{
			UserID:    friend.userID,
			Name:      friend.name,
			HighScore: friend.gameState.Score,
		})
	}

	// Due to the randomization of maps, we'll have to sort it
	sort.Sort(friendByUserID(friends))

	return friends, nil
}

func (db *datastoreSim) DeleteUser(userID string) common.Error {
	db.Lock()
	defer db.Unlock()

	delete(db.Users, userID)
	return nil
}
