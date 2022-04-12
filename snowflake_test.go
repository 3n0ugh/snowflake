package snowflake

import (
	"fmt"
	"strconv"
	"testing"
)

func TestNewNode(t *testing.T) {
	testCases := map[string]struct {
		nodeID       int64
		dataCenterID int64
		err          string
	}{
		"Must_Success": {
			nodeID:       1,
			dataCenterID: 1,
			err:          "",
		},
		"Must_Failure_Negative_NodeID": {
			nodeID:       -1,
			dataCenterID: 1,
			err:          fmt.Sprintf("Node ID must be between 0 and " + strconv.FormatInt(maxNodeID, 10)),
		},
		"Must_Failure_Negative_DataCenterID": {
			nodeID:       1,
			dataCenterID: -1,
			err:          fmt.Sprintf("DataCenter ID must be between 0 and " + strconv.FormatInt(maxDataCenterID, 10)),
		},
	}

	for scenario, tc := range testCases {
		t.Run(scenario, func(t *testing.T) {
			_, err := NewNode(tc.nodeID, tc.dataCenterID)
			if err != nil {
				if tc.err != err.Error() {
					t.Errorf("Expected: %s; Got: %s", tc.err, err.Error())
				}
			}
		})
	}
}

func TestGenerateDuplicateID(t *testing.T) {
	n, err := NewNode(1, 1)
	if err != nil {
		t.Fatalf("Error while creating node: %v", err)
	}

	var holder = make(map[ID]int8)
	for i := 0; i < 1_000_000; i++ {
		id, err := n.Generate()
		if err != nil {
			t.Fatalf("Error while creating id: %v", err)
		}

		if _, exists := holder[id]; exists {
			t.Errorf("id -> %d already exists", id)
		}

		holder[id] = int8(1)
	}
}

// go test -bench=. -count=10 -benchtime=2s .
func BenchmarkGenerateMaxSequence(b *testing.B) {
	node, _ := NewNode(1, 1)

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = node.Generate()
	}
}
