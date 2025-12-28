package sessions_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MuxSuite struct {
	suite.Suite
}

func TestMuxSuite(t *testing.T) {
	suite.Run(t, new(MuxSuite))
}
