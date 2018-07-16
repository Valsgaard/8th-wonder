package datastore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

/**************************************************************************
***************************************************************************
**                                                                       **
**   Testing Suite                                                       **
**                                                                       **
***************************************************************************
**************************************************************************/

type DatastoreTestSuite struct {
	suite.Suite
	ParentCtx context.Context

	Datastore port.Datastore
}

func TestDatastoreSuite(t *testing.T) {
	suite.Run(t, new(DatastoreTestSuite))
}

// Testify/suite workflow
// 1. SetupSuite
// - (Loop for each test)
// 2. SetupTest
// 3. BeforeTest
// 4. TestXXXX
// 5. AfterTest
// 6. TearDownTest
// - (Tests are done)
// 7. TearDownSuite

// Suite setup
func (suite *DatastoreTestSuite) SetupSuite() {
	suite.Datastore = NewDatastoreSimulator()
}

func (suite *DatastoreTestSuite) TearDownSuite() {

}

func (suite *DatastoreTestSuite) SetupTest() {
	suite.prepareSimulatorState()
}

func (suite *DatastoreTestSuite) TearDownTest() {
	// Simulator doesn't require any teardown
}

/**************************************************************************
***************************************************************************
**                                                                       **
**   Testing Data                                                        **
**   Use this data to construct a testing environment in the adapters    **
**                                                                       **
***************************************************************************
**************************************************************************/

var Users = []string{
	"aee6feba-043b-4ba4-a7a4-9d6705595049",
	"bee6feba-043b-4ba4-a7a4-9d6705595049",
	"cee6feba-043b-4ba4-a7a4-9d6705595049",
	"dee6feba-043b-4ba4-a7a4-9d6705595049",
}

var UserNames = []string{
	"bot0",
	"bot1",
	"bot2",
	"bot3",
}

var GameStates = []port.GameState{
	{GamesPlayed: 10, Score: 110},
	{GamesPlayed: 0, Score: 0},
	{GamesPlayed: 0, Score: 0},
	{GamesPlayed: 0, Score: 0},
}

var Friends = [][]string{
	{Users[1], Users[2], Users[3]},
	{},
	{},
	{},
}

/**************************************************************************
***************************************************************************
**                                                                       **
**   Unit Tests                                                          **
**                                                                       **
***************************************************************************
**************************************************************************/

func (suite *DatastoreTestSuite) TestGetUsers() {
	tests := []struct {
		Name            string
		ExpectedSuccess bool
		ExpectedUserIDs []string
		ExpectedNames   []string
	}{
		{
			Name:            "test",
			ExpectedSuccess: true,
			ExpectedUserIDs: []string{
				Users[0],
				Users[1],
				Users[2],
				Users[3],
			},
			ExpectedNames: []string{
				UserNames[0],
				UserNames[1],
				UserNames[2],
				UserNames[3],
			},
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			// Assure test input is correct
			require.Equal(t, len(test.ExpectedNames), len(test.ExpectedUserIDs))

			users, err := suite.Datastore.GetUsers()
			if test.ExpectedSuccess {
				require.Nil(t, err)
				require.NotNil(t, users)

				require.Equal(t, len(test.ExpectedUserIDs), len(users))
				for i, user := range users {
					assert.Equal(t, test.ExpectedUserIDs[i], user.UserID)
					assert.Equal(t, test.ExpectedNames[i], user.Name)
				}
			} else {
				assert.Nil(t, users)
				assert.NotNil(t, err)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}

func (suite *DatastoreTestSuite) TestUserExists() {
	tests := []struct {
		Name            string
		ID              string
		ExpectedSuccess bool
		ExpectedResult  bool
	}{
		{
			Name:            "test",
			ID:              Users[0],
			ExpectedSuccess: true,
			ExpectedResult:  true,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			exists, err := suite.Datastore.UserExists(test.ID)
			if test.ExpectedSuccess {
				require.Nil(t, err)
				assert.Equal(t, test.ExpectedResult, exists)
			} else {
				assert.NotNil(t, err)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}

func (suite *DatastoreTestSuite) TestGetGameState() {
	tests := []struct {
		Name                string
		ID                  string
		ExpectedSuccess     bool
		ExpectedGamesPlayed int
		ExpectedScore       int
	}{
		{
			Name:                "test",
			ID:                  Users[0],
			ExpectedSuccess:     true,
			ExpectedGamesPlayed: GameStates[0].GamesPlayed,
			ExpectedScore:       GameStates[0].Score,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			gameState, err := suite.Datastore.GetGameState(test.ID)
			if test.ExpectedSuccess {
				require.Nil(t, err)
				assert.Equal(t, test.ExpectedGamesPlayed, gameState.GamesPlayed)
				assert.Equal(t, test.ExpectedScore, gameState.Score)
			} else {
				assert.NotNil(t, err)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}

func (suite *DatastoreTestSuite) TestGetFriends() {
	tests := []struct {
		Name            string
		ID              string
		ExpectedSuccess bool
		ExpectedFriends []string
	}{
		{
			Name:            "test",
			ID:              Users[0],
			ExpectedSuccess: true,
			ExpectedFriends: Friends[0],
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			friends, err := suite.Datastore.GetFriends(test.ID)
			if test.ExpectedSuccess {
				require.Nil(t, err)
				require.NotNil(t, friends)

				require.Equal(t, len(test.ExpectedFriends), len(friends))
				for i, friend := range friends {
					assert.Equal(t, test.ExpectedFriends[i], friend.UserID)
				}
			} else {
				assert.Nil(t, friends)
				assert.NotNil(t, err)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}

func (suite *DatastoreTestSuite) TestNewUser() {
	tests := []struct {
		Name            string
		ID              string
		UserName        string
		ExpectedSuccess bool
	}{
		{
			Name:            "test",
			ID:              "fee6feba-043b-4ba4-a7a4-9d6705595049",
			UserName:        "flaf",
			ExpectedSuccess: true,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			user, err := suite.Datastore.NewUser(test.ID, test.UserName)

			if test.ExpectedSuccess {
				assert.Nil(t, err)
				if assert.NotNil(t, user) {
					assert.Equal(t, test.ID, user.UserID)
					assert.Equal(t, test.UserName, user.Name)
				}

				users, err := suite.Datastore.GetUsers()
				require.Nil(t, err)

				require.Equal(t, len(Users)+1, len(users))
				assert.Equal(t, test.ID, users[len(users)-1].UserID)
				assert.Equal(t, test.UserName, users[len(users)-1].Name)
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, user)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}

func (suite *DatastoreTestSuite) TestUpdateGameState() {
	tests := []struct {
		Name            string
		ID              string
		GamesPlayed     int
		Score           int
		ExpectedSuccess bool
	}{
		{
			Name:            "test",
			ID:              Users[3],
			GamesPlayed:     1,
			Score:           200,
			ExpectedSuccess: true,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			err := suite.Datastore.UpdateGameState(test.ID, test.GamesPlayed, test.Score)
			if test.ExpectedSuccess {
				require.Nil(t, err)

				// Read to verify changes
				gameState, err := suite.Datastore.GetGameState(test.ID)
				require.Nil(t, err)

				assert.Equal(t, test.GamesPlayed, gameState.GamesPlayed)
				assert.Equal(t, test.Score, gameState.Score)

			} else {
				assert.NotNil(t, err)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}

func (suite *DatastoreTestSuite) TestUpdateFriends() {
	tests := []struct {
		Name            string
		ID              string
		Friends         []string
		ExpectedSuccess bool
	}{
		{
			Name: "test",
			ID:   Users[3],
			Friends: []string{
				Users[0],
				Users[1],
			},
			ExpectedSuccess: true,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			err := suite.Datastore.UpdateFriends(test.ID, test.Friends)
			if test.ExpectedSuccess {
				require.Nil(t, err)

				friends, err := suite.Datastore.GetFriends(test.ID)

				require.Nil(t, err)
				require.NotNil(t, friends)

				require.Equal(t, len(test.Friends), len(friends))
				for i, friend := range friends {
					assert.Equal(t, test.Friends[i], friend.UserID)
				}
			} else {
				assert.NotNil(t, err)
			}
		}

		suite.T().Run(test.Name, fn)
	}

}

func (suite *DatastoreTestSuite) TestDeleteUser() {
	tests := []struct {
		Name            string
		ID              string
		ExpectedSuccess bool
	}{
		{
			Name:            "test",
			ID:              Users[0],
			ExpectedSuccess: true,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			err := suite.Datastore.DeleteUser(test.ID)
			if test.ExpectedSuccess {
				require.Nil(t, err)

				users, err := suite.Datastore.GetUsers()
				require.Nil(t, err)
				require.NotNil(t, users)

				assert.Equal(t, len(Users)-1, len(users))
				assert.Equal(t, Users[1], users[0].UserID)
			} else {
				assert.NotNil(t, err)
			}
		}

		suite.T().Run(test.Name, fn)
	}
}
