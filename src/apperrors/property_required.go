package apperrors

import "fmt"

type PropertyRequiredError struct {
	Entity   string
	Property string
}

func NewPropertyRequiredError(entity, property string) PropertyRequiredError {
	return PropertyRequiredError{
		Entity:   entity,
		Property: property,
	}
}

func (e PropertyRequiredError) Error() string {
	return fmt.Sprintf("'%s' property is required for %s entity", e.Property, e.Entity)
}
