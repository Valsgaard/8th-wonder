// Note: Endpoints test are called in isolated package as it's accessing a
//       library from an outer layer (Datastore).

package endpoints_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/valsgaard/interview-case/backend/common"
	"github.com/valsgaard/interview-case/backend/datastore"
	"github.com/valsgaard/interview-case/backend/endpoints/port"
)

type EndpointsTestSuite struct {
	suite.Suite
	ParentCtx context.Context

	Users []string
}

func TestEndpointsSuite(t *testing.T) {
	suite.Run(t, new(EndpointsTestSuite))
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
func (suite *EndpointsTestSuite) SetupSuite() {
	suite.Users = []string{
		"aee6feba-043b-4ba4-a7a4-9d6705595049",
		"bee6feba-043b-4ba4-a7a4-9d6705595049",
		"cee6feba-043b-4ba4-a7a4-9d6705595049",
		"dee6feba-043b-4ba4-a7a4-9d6705595049",
	}
}

func (suite *EndpointsTestSuite) TearDownSuite() {

}

func (suite *EndpointsTestSuite) SetupTest() {
	datastore := datastore.NewDatastoreSimulator()

	ctx := context.Background()
	ctx = port.SetDatastore(ctx, datastore)
	ctx = common.SetLog(ctx, common.NewLog("test", os.Stdout))
	suite.ParentCtx = ctx

	for i := range suite.Users {
		datastore.NewUser(suite.Users[i], fmt.Sprintf("bot%d", i))
	}

	// Set updated friends
	port.GetDatastore(suite.ParentCtx).UpdateFriends(
		suite.Users[0],
		[]string{
			suite.Users[1],
			suite.Users[2],
			suite.Users[3],
		},
	)

	// Set game state
	port.GetDatastore(suite.ParentCtx).UpdateGameState(
		suite.Users[0],
		10,
		110,
	)
}

func (suite *EndpointsTestSuite) TearDownTest() {
	// Remove users
	for i := range suite.Users {
		port.GetDatastore(suite.ParentCtx).DeleteUser(suite.Users[i])
	}
}
