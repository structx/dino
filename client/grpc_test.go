package client_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/structx/dino/client"
)

type ClientSuite struct {
	suite.Suite

	cli client.Client
}

func (suite *ClientSuite) SetupSuite() {
	_ = suite.cli
}

func (suite *ClientSuite) TestTunnelAdd() {}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
