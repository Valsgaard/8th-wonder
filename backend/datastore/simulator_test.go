package datastore

// Preare the state of the Simulator adapter for each test itereration
// Note that the preperation overwrites everything, and requires no teardown
func (suite *DatastoreTestSuite) prepareSimulatorState() {
	db := suite.Datastore.(*datastoreSim)
	db.Users = map[string]*datastoreUser{
		Users[0]: &datastoreUser{
			userID:    Users[0],
			name:      UserNames[0],
			gameState: GameStates[0],
			friendIDs: []string{
				Users[1],
				Users[2],
				Users[3],
			},
		},

		Users[1]: &datastoreUser{
			userID:    Users[1],
			name:      UserNames[1],
			gameState: GameStates[1],
			friendIDs: []string{},
		},

		Users[2]: &datastoreUser{
			userID:    Users[2],
			name:      UserNames[2],
			gameState: GameStates[2],
			friendIDs: []string{},
		},

		Users[3]: &datastoreUser{
			userID:    Users[3],
			name:      UserNames[3],
			gameState: GameStates[3],
			friendIDs: []string{},
		},
	}
}
