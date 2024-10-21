package validator

import "server/internal/repository"

func ValidateTrip(v *Validator, trip *repository.InsertTripParams) {
	v.Check(trip.Name != "", "name", "must be provided")
	v.Check(len(trip.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(trip.Description.String != "", "description", "must be provided")
	v.Check(len(trip.Description.String) <= 500, "description", "must not be more than 500 bytes long")

	_, hasError := v.Errors["description"]
	if !hasError {
		trip.Description.Valid = true;
	}
}
