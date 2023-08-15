package apperrors

import "fmt"

type NoShipFoundErr struct {
	id int32
}

func (e *NoShipFoundErr) Error() string {
	return fmt.Sprintf("no matching ship found for id: %d", e.id)
}

func NewNoShipFoundErr(id int32) *NoShipFoundErr {
	return &NoShipFoundErr{id: id}
}
