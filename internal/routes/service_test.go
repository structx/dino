package routes

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RouteServiceSuite struct {
	suite.Suite
}

func (suite *RouteServiceSuite) SetupSuite() {}

func TestRouteServiceSuite(t *testing.T) {
	suite.Run(t, new(RouteServiceSuite))
}
