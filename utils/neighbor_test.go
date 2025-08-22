package utils_test

import (
	"fmt"
	"github.com/matrix-go/bitcoin/core"
	"github.com/matrix-go/bitcoin/utils"
	"testing"
)

func TestFindNeighbors(t *testing.T) {
	neighbors := utils.FindNeighbors(
		"127.0.0.1", 5000,
		core.KNeighborIPRangeStart, core.KNeighborIPRangeEnd,
		core.KBlockchainPortStart, core.KBlockchainPortEnd)
	for _, neighbor := range neighbors {
		fmt.Printf("neighbor: %v\n", neighbor)
	}
}
