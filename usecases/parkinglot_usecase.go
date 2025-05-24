package usecases

import (
	"errors"
	"fmt"
	"strings"
	"submit_do_it/constants"
	"submit_do_it/domain"
)

type parkinglotUsecaseImpl struct {
	pl *domain.ParkingLot
}

type ParkinglotUsecase interface {
	Park(vehicleType constants.VehicleType, vehicleNumber string) (string, error)
	Unpark(spotId, vehicleNumber string) error
	AvailableSpot(vehicleType constants.VehicleType) int
	SearchVehicle(vehicleNumber string) (string, error)
}

func NewParkingLotUsecase(floors, rows, columns int, layoutTemplate [][]string) ParkinglotUsecase {
	lot := &domain.ParkingLot{
		Floors:         floors,
		Rows:           rows,
		Columns:        columns,
		Layout:         make([][][]*domain.Spot, floors),
		VehicleMap:     make(map[string]string),
		LastSpotMap:    make(map[string]string),
		AvailableSpots: make(map[constants.VehicleType]map[string]*domain.Spot),
	}

	for _, vt := range []constants.VehicleType{
		constants.Bicycle,
		constants.Motorcycle,
		constants.Automobile,
	} {
		lot.AvailableSpots[vt] = make(map[string]*domain.Spot)
	}

	for f := 0; f < floors; f++ {
		lot.Layout[f] = make([][]*domain.Spot, rows)
		for r := 0; r < rows; r++ {
			lot.Layout[f][r] = make([]*domain.Spot, columns)
			for c := 0; c < columns; c++ {
				code := layoutTemplate[r][c]
				parts := strings.Split(code, "-")
				spotType := parts[0]
				active := parts[1] == "1"

				var vt constants.VehicleType
				switch spotType {
				case "B":
					vt = constants.Bicycle
				case "M":
					vt = constants.Motorcycle
				case "A":
					vt = constants.Automobile
				default:
					active = false
				}

				spot := &domain.Spot{
					Floor:    f,
					Row:      r,
					Col:      c,
					SpotType: vt,
					Active:   active,
				}

				lot.Layout[f][r][c] = spot

				if active {
					spotID := spot.ID()
					lot.AvailableSpots[vt][spotID] = spot
				}
			}
		}
	}
	return &parkinglotUsecaseImpl{
		pl: lot,
	}
}

func (pu *parkinglotUsecaseImpl) Park(vehicleType constants.VehicleType, vehicleNumber string) (string, error) {
	pu.pl.Mutx.Lock()
	defer pu.pl.Mutx.Unlock()

	if _, exists := pu.pl.VehicleMap[vehicleNumber]; exists {
		return "", errors.New("vehicle already parked")
	}

	for spotID, spot := range pu.pl.AvailableSpots[vehicleType] {
		spot.Mutx.Lock()
		if spot.Occupied {
			spot.Mutx.Unlock()
			continue
		}
		spot.Occupied = true
		spot.VehicleNumber = vehicleNumber
		spot.Mutx.Unlock()

		pu.pl.VehicleMap[vehicleNumber] = spotID
		pu.pl.LastSpotMap[vehicleNumber] = spotID
		delete(pu.pl.AvailableSpots[vehicleType], spotID)
		return spotID, nil
	}

	return "", errors.New("no available parking spot for vehicle type")
}

func (pu *parkinglotUsecaseImpl) Unpark(spotID, vehicleNumber string) error {
	pu.pl.Mutx.Lock()
	defer pu.pl.Mutx.Unlock()

	if pu.pl.VehicleMap[vehicleNumber] != spotID {
		return errors.New("vehicle not found at specified spot")
	}

	_ = strings.Split(spotID, "-")
	var f, r, c int
	fmt.Sscanf(spotID, "%d-%d-%d", &f, &r, &c)

	spot := pu.pl.Layout[f][r][c]
	spot.Mutx.Lock()
	defer spot.Mutx.Unlock()

	if !spot.Occupied || spot.VehicleNumber != vehicleNumber {
		return errors.New("spot not occupied by this vehicle")
	}

	spot.Occupied = false
	spot.VehicleNumber = ""
	delete(pu.pl.VehicleMap, vehicleNumber)
	pu.pl.AvailableSpots[spot.SpotType][spotID] = spot

	return nil
}

func (pu *parkinglotUsecaseImpl) AvailableSpot(vehicleType constants.VehicleType) int {
	pu.pl.Mutx.RLock()
	defer pu.pl.Mutx.RUnlock()
	return len(pu.pl.AvailableSpots[vehicleType])
}

func (pu *parkinglotUsecaseImpl) SearchVehicle(vehicleNumber string) (string, error) {
	pu.pl.Mutx.RLock()
	defer pu.pl.Mutx.RUnlock()

	if spot, ok := pu.pl.VehicleMap[vehicleNumber]; ok {
		return spot, nil
	}
	if last, ok := pu.pl.LastSpotMap[vehicleNumber]; ok {
		return last, nil
	}
	return "", errors.New("vehicle not found")
}
