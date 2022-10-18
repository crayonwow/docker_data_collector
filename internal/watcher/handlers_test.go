package watcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_prepareStatTemplate(t *testing.T) {
	se := statsEntry{
		Container:        "1",
		Name:             "1",
		ID:               "1",
		CPUPercentage:    1,
		Memory:           1,
		MemoryLimit:      1,
		MemoryPercentage: 1,
		NetworkRx:        1,
		NetworkTx:        1,
		NetworkIO:        "1",
		BlockRead:        1,
		BlockWrite:       1,
		PidsCurrent:      1,
		IsInvalid:        false,
	}
	message, err := prepareStatTemplate(se)
	require.NoError(t, err)
	require.Equal(t, `Name: 1
CPU %:    1.00%
MEM %:    1.00%
NETWORK : 1
`, message)
}
