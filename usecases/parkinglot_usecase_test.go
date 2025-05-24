package usecases

import (
	"fmt"
	"testing"

	"submit_do_it/constants"
)

func TestNewParkingLotUsecase_BasicLayout(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1", "M-1", "A-1"},
		{"B-0", "M-1", "A-0"},
	}
	u := NewParkingLotUsecase(1, 2, 3, layoutTemplate)
	impl, ok := u.(*parkinglotUsecaseImpl)
	if !ok {
		t.Fatalf("expected *parkinglotUsecaseImpl, got %T", u)
	}
	pl := impl.pl

	if pl.Floors != 1 || pl.Rows != 2 || pl.Columns != 3 {
		t.Errorf("unexpected lot dimensions: %+v", pl)
	}

	// Check spot types and activeness
	tests := []struct {
		f, r, c  int
		vt       constants.VehicleType
		active   bool
		template string
	}{
		{0, 0, 0, constants.Bicycle, true, "B-1"},
		{0, 0, 1, constants.Motorcycle, true, "M-1"},
		{0, 0, 2, constants.Automobile, true, "A-1"},
		{0, 1, 0, constants.Bicycle, false, "B-0"},
		{0, 1, 1, constants.Motorcycle, true, "M-1"},
		{0, 1, 2, constants.Automobile, false, "A-0"},
	}
	for _, tt := range tests {
		spot := pl.Layout[tt.f][tt.r][tt.c]
		if spot.SpotType != tt.vt {
			t.Errorf("spot at %d-%d-%d: expected type %v, got %v", tt.f, tt.r, tt.c, tt.vt, spot.SpotType)
		}
		if spot.Active != tt.active {
			t.Errorf("spot at %d-%d-%d: expected active %v, got %v", tt.f, tt.r, tt.c, tt.active, spot.Active)
		}
	}

	// Check AvailableSpots map
	expectedAvailable := map[constants.VehicleType]int{
		constants.Bicycle:    1,
		constants.Motorcycle: 2,
		constants.Automobile: 1,
	}
	for vt, want := range expectedAvailable {
		got := len(pl.AvailableSpots[vt])
		if got != want {
			t.Errorf("AvailableSpots[%v]: want %d, got %d", vt, want, got)
		}
	}
}

func TestNewParkingLotUsecase_MultipleFloors(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1", "M-1"},
		{"A-1", "B-0"},
	}
	u := NewParkingLotUsecase(2, 2, 2, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)
	pl := impl.pl

	if len(pl.Layout) != 2 {
		t.Errorf("expected 2 floors, got %d", len(pl.Layout))
	}
	for f := 0; f < 2; f++ {
		for r := 0; r < 2; r++ {
			for c := 0; c < 2; c++ {
				spot := pl.Layout[f][r][c]
				if spot.Floor != f || spot.Row != r || spot.Col != c {
					t.Errorf("spot at [%d][%d][%d] has wrong indices: %+v", f, r, c, spot)
				}
			}
		}
	}
}

func TestNewParkingLotUsecase_UnknownSpotType(t *testing.T) {
	layoutTemplate := [][]string{
		{"X-1", "B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 2, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)
	pl := impl.pl

	spot := pl.Layout[0][0][0]
	if spot.Active {
		t.Errorf("unknown spot type should not be active")
	}
	// if spot.SpotType != constants.VehicleType(0) {
	// 	t.Errorf("unknown spot type should be zero value, got %v", spot.SpotType)
	// }
}

func TestNewParkingLotUsecase_EmptyLayout(t *testing.T) {
	layoutTemplate := [][]string{}
	u := NewParkingLotUsecase(1, 0, 0, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)
	pl := impl.pl

	if len(pl.Layout) != 1 {
		t.Errorf("expected 1 floor, got %d", len(pl.Layout))
	}
	if len(pl.Layout[0]) != 0 {
		t.Errorf("expected 0 rows, got %d", len(pl.Layout[0]))
	}
	for vt := range pl.AvailableSpots {
		if len(pl.AvailableSpots[vt]) != 0 {
			t.Errorf("expected 0 available spots for %v", vt)
		}
	}
}

func TestNewParkingLotUsecase_MapsInitialized(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)
	pl := impl.pl

	if pl.VehicleMap == nil || pl.LastSpotMap == nil {
		t.Errorf("VehicleMap or LastSpotMap not initialized")
	}
	for _, vt := range []constants.VehicleType{constants.Bicycle, constants.Motorcycle, constants.Automobile} {
		if pl.AvailableSpots[vt] == nil {
			t.Errorf("AvailableSpots[%v] not initialized", vt)
		}
	}
}

func TestNewParkingLotUsecase_SpotIDs(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1", "M-1"},
	}
	u := NewParkingLotUsecase(1, 1, 2, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)
	pl := impl.pl

	for vt, spots := range pl.AvailableSpots {
		for spotID, spot := range spots {
			expectedID := spot.ID()
			if spotID != expectedID {
				t.Errorf("spotID mismatch: got %s, want %s", spotID, expectedID)
			}
			if spot.SpotType != vt {
				t.Errorf("spot type mismatch: got %v, want %v", spot.SpotType, vt)
			}
		}
	}
}

func TestParkinglotUsecaseImpl_Park_Success(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1", "M-1", "A-1"},
	}
	u := NewParkingLotUsecase(1, 1, 3, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	spotID, err := impl.Park(constants.Bicycle, "BIKE123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if spotID == "" {
		t.Errorf("expected a spotID, got empty string")
	}
	// Check that the vehicle is mapped
	if got, ok := impl.pl.VehicleMap["BIKE123"]; !ok || got != spotID {
		t.Errorf("vehicle not mapped correctly, got %v, want %v", got, spotID)
	}
	// Check that the spot is now occupied
	spot := impl.pl.Layout[0][0][0]
	if !spot.Occupied || spot.VehicleNumber != "BIKE123" {
		t.Errorf("spot not marked as occupied or wrong vehicle number")
	}
	// Check that the spot is removed from AvailableSpots
	if _, ok := impl.pl.AvailableSpots[constants.Bicycle][spotID]; ok {
		t.Errorf("spot should be removed from AvailableSpots after parking")
	}
}

func TestParkinglotUsecaseImpl_Park_AlreadyParked(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	spotID, err := impl.Park(constants.Bicycle, "BIKE123")
	if err != nil {
		t.Fatalf("first park should succeed, got %v", err)
	}
	spotID2, err2 := impl.Park(constants.Bicycle, "BIKE123")
	if err2 == nil {
		t.Errorf("expected error for already parked vehicle, got nil")
	}
	if spotID2 != "" {
		t.Errorf("expected empty spotID for already parked vehicle, got %v", spotID2)
	}
	// Should still be mapped to the first spot
	if got := impl.pl.VehicleMap["BIKE123"]; got != spotID {
		t.Errorf("vehicle map changed unexpectedly, got %v, want %v", got, spotID)
	}
}

func TestParkinglotUsecaseImpl_Park_NoAvailableSpot(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	_, err := impl.Park(constants.Bicycle, "BIKE123")
	if err != nil {
		t.Fatalf("first park should succeed, got %v", err)
	}
	_, err2 := impl.Park(constants.Bicycle, "BIKE456")
	if err2 == nil {
		t.Errorf("expected error when no available spot, got nil")
	}
}

func TestParkinglotUsecaseImpl_Park_WrongVehicleType(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	_, err := impl.Park(constants.Automobile, "CAR123")
	fmt.Println(err)
	if err == nil {
		t.Errorf("expected error when no spot for vehicle type, got nil")
	}
}
func TestParkinglotUsecaseImpl_Unpark_Success(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1", "M-1"},
	}
	u := NewParkingLotUsecase(1, 1, 2, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	// Park a bicycle
	spotID, err := impl.Park(constants.Bicycle, "BIKE123")
	if err != nil {
		t.Fatalf("Park failed: %v", err)
	}

	// Unpark the bicycle
	err = impl.Unpark(spotID, "BIKE123")
	if err != nil {
		t.Errorf("expected no error on unpark, got %v", err)
	}

	// Spot should be available again
	if _, ok := impl.pl.AvailableSpots[constants.Bicycle][spotID]; !ok {
		t.Errorf("spot should be available after unpark")
	}
	// Vehicle should be removed from VehicleMap
	if _, ok := impl.pl.VehicleMap["BIKE123"]; ok {
		t.Errorf("vehicle should be removed from VehicleMap after unpark")
	}
	// Spot should not be occupied
	spot := impl.pl.Layout[0][0][0]
	if spot.Occupied || spot.VehicleNumber != "" {
		t.Errorf("spot should not be occupied after unpark")
	}
}

func TestParkinglotUsecaseImpl_Unpark_WrongSpotID(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	_, err := impl.Park(constants.Bicycle, "BIKE123")
	if err != nil {
		t.Fatalf("Park failed: %v", err)
	}

	// Try to unpark with wrong spotID
	err = impl.Unpark("0-0-99", "BIKE123")
	if err == nil {
		t.Errorf("expected error for wrong spotID, got nil")
	}
	// Should still be parked
	if _, ok := impl.pl.VehicleMap["BIKE123"]; !ok {
		t.Errorf("vehicle should still be parked")
	}
}

func TestParkinglotUsecaseImpl_Unpark_WrongVehicleNumber(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	spotID, err := impl.Park(constants.Bicycle, "BIKE123")
	if err != nil {
		t.Fatalf("Park failed: %v", err)
	}

	// Try to unpark with wrong vehicle number
	err = impl.Unpark(spotID, "BIKE999")
	if err == nil {
		t.Errorf("expected error for wrong vehicle number, got nil")
	}
	// Should still be parked
	if _, ok := impl.pl.VehicleMap["BIKE123"]; !ok {
		t.Errorf("vehicle should still be parked")
	}
}

func TestParkinglotUsecaseImpl_Unpark_SpotNotOccupied(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	// Try to unpark when nothing is parked
	err := impl.Unpark("0-0-0", "BIKE123")
	if err == nil {
		t.Errorf("expected error when spot not occupied, got nil")
	}
}

func TestParkinglotUsecaseImpl_Unpark_SpotOccupiedByOtherVehicle(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	spotID, err := impl.Park(constants.Bicycle, "BIKE123")
	if err != nil {
		t.Fatalf("Park failed: %v", err)
	}

	// Manually change the vehicle number at the spot
	spot := impl.pl.Layout[0][0][0]
	spot.VehicleNumber = "BIKE999"

	// Try to unpark with original vehicle number
	err = impl.Unpark(spotID, "BIKE123")
	if err == nil {
		t.Errorf("expected error when spot occupied by another vehicle, got nil")
	}
}
func TestParkinglotUsecaseImpl_AvailableSpot_Basic(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1", "M-1", "A-1"},
		{"B-1", "M-0", "A-1"},
	}
	u := NewParkingLotUsecase(1, 2, 3, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	tests := []struct {
		vt     constants.VehicleType
		expect int
	}{
		{constants.Bicycle, 2},
		{constants.Motorcycle, 1},
		{constants.Automobile, 2},
	}
	for _, tt := range tests {
		got := impl.AvailableSpot(tt.vt)
		if got != tt.expect {
			t.Errorf("AvailableSpot(%v): got %d, want %d", tt.vt, got, tt.expect)
		}
	}
}

func TestParkinglotUsecaseImpl_AvailableSpot_AfterParkingAndUnparking(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1", "B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 2, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	// Initially 2 bicycle spots
	if got := impl.AvailableSpot(constants.Bicycle); got != 2 {
		t.Errorf("expected 2 available, got %d", got)
	}

	spotID, err := impl.Park(constants.Bicycle, "BIKE1")
	if err != nil {
		t.Fatalf("Park failed: %v", err)
	}
	if got := impl.AvailableSpot(constants.Bicycle); got != 1 {
		t.Errorf("expected 1 available after parking, got %d", got)
	}

	err = impl.Unpark(spotID, "BIKE1")
	if err != nil {
		t.Fatalf("Unpark failed: %v", err)
	}
	if got := impl.AvailableSpot(constants.Bicycle); got != 2 {
		t.Errorf("expected 2 available after unparking, got %d", got)
	}
}

func TestParkinglotUsecaseImpl_AvailableSpot_EmptyLot(t *testing.T) {
	layoutTemplate := [][]string{}
	u := NewParkingLotUsecase(1, 0, 0, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	for _, vt := range []constants.VehicleType{constants.Bicycle, constants.Motorcycle, constants.Automobile} {
		if got := impl.AvailableSpot(vt); got != 0 {
			t.Errorf("expected 0 available for %v, got %d", vt, got)
		}
	}
}
func TestParkinglotUsecaseImpl_SearchVehicle_Parked(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1", "M-1"},
	}
	u := NewParkingLotUsecase(1, 1, 2, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	spotID, err := impl.Park(constants.Bicycle, "BIKE123")
	if err != nil {
		t.Fatalf("Park failed: %v", err)
	}

	foundSpot, err := impl.SearchVehicle("BIKE123")
	if err != nil {
		t.Errorf("expected to find vehicle, got error: %v", err)
	}
	if foundSpot != spotID {
		t.Errorf("expected spotID %s, got %s", spotID, foundSpot)
	}
}

func TestParkinglotUsecaseImpl_SearchVehicle_LastSpot(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	spotID, err := impl.Park(constants.Bicycle, "BIKE123")
	if err != nil {
		t.Fatalf("Park failed: %v", err)
	}
	err = impl.Unpark(spotID, "BIKE123")
	if err != nil {
		t.Fatalf("Unpark failed: %v", err)
	}

	foundSpot, err := impl.SearchVehicle("BIKE123")
	if err != nil {
		t.Errorf("expected to find last spot, got error: %v", err)
	}
	if foundSpot != spotID {
		t.Errorf("expected last spotID %s, got %s", spotID, foundSpot)
	}
}

func TestParkinglotUsecaseImpl_SearchVehicle_NotFound(t *testing.T) {
	layoutTemplate := [][]string{
		{"B-1"},
	}
	u := NewParkingLotUsecase(1, 1, 1, layoutTemplate)
	impl := u.(*parkinglotUsecaseImpl)

	_, err := impl.SearchVehicle("UNKNOWN123")
	if err == nil {
		t.Errorf("expected error for unknown vehicle, got nil")
	}
}
