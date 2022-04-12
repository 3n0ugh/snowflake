package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

const (
	nodeIDBits       = uint8(5)  // 5bit nodeID
	dataCenterIDBits = uint8(5)  // 5bit dataCenterID
	sequenceBits     = uint8(12) // 12bit sequence

	// The maximum value of the node ID used to prevent overflow
	maxNodeID = int64(-1) ^ (int64(-1) << nodeIDBits)
	// The maximum value of the data center ID used to prevent overflow
	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	// The maximum value of sequence used to prevent overflow
	maxSequence = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft = uint8(22) // timeLeft = nodeIDBits + sequenceBits // Timestamp offset to the left
	dataLeft = uint8(17) // dataLeft = dataCenterIDBits + sequenceBits
	nodeLeft = uint8(12) // nodeLeft = sequenceBits // Node IDx offset to the left
	// 2020-05-20 08:00:00 +0800 CST
	epoch = int64(1589923200000) // Constant timestamp (ms)
)

type Node struct {
	mu sync.Mutex

	// Record the time stamp of the last ID
	LastStamp int64

	// ID of the node
	NodeID int64

	// DataCenterID of the node
	DataCenterID int64

	// ID sequence numbers that have been generated in the current millisecond
	// (accumulated from 0)A maximum of 4096 IDs are generated within 1 millisecond
	Sequence int64
}

type ID uint64

func NewNode(nodeID, dataCenterID int64) (*Node, error) {
	if nodeID < 0 || nodeID > maxNodeID {
		return nil, errors.New("Node ID must be between 0 and " + strconv.FormatInt(maxNodeID, 10))
	}

	if dataCenterID < 0 || dataCenterID > maxDataCenterID {
		return nil, errors.New("DataCenter ID must be between 0 and " + strconv.FormatInt(maxDataCenterID, 10))
	}

	return &Node{
		NodeID:       nodeID,
		LastStamp:    0,
		Sequence:     0,
		DataCenterID: dataCenterID,
	}, nil
}

func (w *Node) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (w *Node) Generate() (ID, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.generate()
}

func (w *Node) generate() (ID, error) {
	timeStamp := w.getMilliSeconds()
	if timeStamp < w.LastStamp {
		return 0, errors.New("time is moving backwards,waiting until")
	}

	if w.LastStamp == timeStamp {
		w.Sequence = (w.Sequence + 1) & maxSequence

		if w.Sequence == 0 {
			for timeStamp <= w.LastStamp {
				timeStamp = w.getMilliSeconds()
			}
		}
	} else {
		w.Sequence = 0
	}

	w.LastStamp = timeStamp

	id := ID(((timeStamp - epoch) << timeLeft) |
		(w.DataCenterID << dataLeft) |
		(w.NodeID << nodeLeft) |
		w.Sequence)

	return id, nil
}

// String returns a string of the snowflake ID
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// UInt64 returns a uint64 of the snowflake ID
func (f ID) UInt64() uint64 {
	return uint64(f)
}

func DecomposeID(id ID) map[string]uint64 {
	const maskSequence = uint64((1<<dataLeft - 1) >> sequenceBits)
	const maskMachineID = uint64(1<<dataLeft - 1)

	t := id.UInt64() >> timeLeft
	sequence := id.UInt64() & maskSequence
	machineID := id.UInt64() & maskMachineID >> sequenceBits
	return map[string]uint64{
		"id":         id.UInt64(),
		"time":       t,
		"sequence":   sequence,
		"machine-id": machineID,
	}
}
