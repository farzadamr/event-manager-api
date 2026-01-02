package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func DateValidator(fld validator.FieldLevel) bool {
	date, ok := fld.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	now := time.Now().UTC()
	sixMonthLater := now.AddDate(0, 6, 0)

	return !date.Before(now) && !date.After(sixMonthLater)
}
