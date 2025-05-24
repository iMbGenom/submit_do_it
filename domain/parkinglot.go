package domain

import (
	"fmt"
	"submit_do_it/constants"
	"sync"
)

type Spot struct {
	Floor    int
	Row      int
	Col      int
	SpotType constants.VehicleType
	Active   bool

	VehicleNumber string
	Occupied      bool

	Mutx sync.Mutex
}

type ParkingLot struct {
	Floors  int
	Rows    int
	Columns int
	Layout  [][][]*Spot

	VehicleMap     map[string]string
	LastSpotMap    map[string]string
	AvailableSpots map[constants.VehicleType]map[string]*Spot

	Mutx sync.RWMutex // protects vehicleMap, lastSpotMap, availableSpots
}

// Returns spot ID as "floor-row-col"
func (s *Spot) ID() string {
	return fmt.Sprintf("%d-%d-%d", s.Floor, s.Row, s.Col)
}
