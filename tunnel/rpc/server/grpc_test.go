package server

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ReverseTunnelSuite struct {
	suite.Suite
}

func (suite *ReverseTunnelSuite) SetupSuite() {}

func TestGRPCSuite(t *testing.T) {
	suite.Run(t, new(ReverseTunnelSuite))
}
