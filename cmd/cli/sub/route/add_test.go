package route

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"soft.structx.io/dino/cmd/cli/sub"
)

func TestAddCmd(t *testing.T) {
	assert := assert.New(t)
	cmd := routeCmd

	tt := []struct {
		args     []string
		err      error
		expected string
	}{
		{
			args:     []string{"route", "add"},
			err:      nil,
			expected: "",
		},
	}

	for _, tc := range tt {
		actual, err := sub.CmdExecute(t, cmd, tc.args)
		assert.Equal(tc.err, err)

		assert.Equal(tc.expected, actual)
	}
}
