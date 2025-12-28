package client

import (
	"fmt"
	"strconv"
)

const (
	maxPort = 65535
	minPort = 1

	base = 10
)

func validateAndParsePort(port uint32) (string, error) {
	if port < minPort || port > maxPort {
		return "", fmt.Errorf("port number %d is outside range %d-%d", port, minPort, maxPort)
	}

	portStr := strconv.FormatUint(uint64(port), base)
	return portStr, nil
}
