package main

import (
	"fmt"
	"submit_do_it/constants"
	"submit_do_it/usecases"
)

type vehicleInfo struct {
	vehicleType   constants.VehicleType
	vehicleNumber string
}

func main() {
	layout := [][]string{
		{"B-1", "M-1", "A-1"},
		{"B-1", "X-0", "A-1"},
		{"M-1", "A-1", "B-1"},
	}
	vehicles := []vehicleInfo{
		{vehicleType: constants.Bicycle, vehicleNumber: "BIKE123"},
		{vehicleType: constants.Motorcycle, vehicleNumber: "MOTO123"},
		{vehicleType: constants.Automobile, vehicleNumber: "CAR123"},
	}

	for _, v := range vehicles {
		runVehicle(layout, v.vehicleType, v.vehicleNumber)
	}
}

func runVehicle(layout [][]string, vehicleType constants.VehicleType, vehicleNumber string) {
	pl := usecases.NewParkingLotUsecase(2, 3, 3, layout)

	fmt.Printf("Vehicle number: %#v\n", vehicleNumber)

	// Park vehicle. Given a vehicle type, assign an empty parking spot id and map the vehicleNumber. spotId is floor-row-column. If no free spot is found, return an error.
	spotId, err := pl.Park(vehicleType, vehicleNumber)
	if err != nil {
		fmt.Printf("Park failed: %#v\n", err)
	}
	fmt.Printf("%s Parked at: %#v\n", vehicleType, spotId)

	// Unpark vehicle. Removes vehicle from parking spot. Return an error for failure to unpark a vehicle.
	err = pl.Unpark(spotId, vehicleNumber)
	if err != nil {
		fmt.Printf("Unpark failed: %#v\n", err)
	}

	// Available spot. Display the free spots for each vehicle type.
	fmt.Printf("Available %s spots: %#v\n", vehicleType, pl.AvailableSpot(vehicleType))

	// Search vehicle. If the vehicle has been unparked, get its last spotId
	spot, err := pl.SearchVehicle(vehicleNumber)
	if err != nil {
		fmt.Printf("SearchVehicle failed: %#v\n", err)
	}
	fmt.Printf("Vehicle location: %#v\n\n", spot)
}
