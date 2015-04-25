package dataTypes

import (
	"sync"
	"time"
)

const (
	N_FLOORS        = 4
	N_BUTTONS       = 4
	//NUMBER_OF_ELEVS = 2	// changed

	// Global timeout const
	TIMEOUT = 1 * time.Second
)

var (
	NumberOfElevatorsConnected 	int	// changed
	ElevatorID         		int	// changed
)

type (
	Matrix [][]int
	Array []int
	ElevButtonTypeT int

	// Map for storing addresses of peers in group
	PeerMap struct {
		Mu sync.Mutex
		M  map[int]time.Time
	}

	// Struct for sending data over network
	Message struct {	//changed 
		Head     string
		Order    []int
		Table    [][]int
		Cost     int
		ID       int
		WhichExternalPanelPressed int
		T        time.Time

		// added
		Word	string
	}
)

func NewExternalQueue() [][]int {	// changed
	t := make([][]int, N_FLOORS)
	for i := range t {
		t[i] = make([]int, 2)
		for j := range t[i] {
			t[i][j] = 0
		}
	}
	return t
}

func NewInternalQueue() []int {
	t := make([]int, N_FLOORS)
	return t
}
