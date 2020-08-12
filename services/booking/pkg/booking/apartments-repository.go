package booking

import (
	"context"
)

type apartmentsRepository struct {
}

func (a *apartmentsRepository) GetApartmentById(ctx context.Context, apartmentId string) (*Apartment, error) {
	panic("implement me")
}
