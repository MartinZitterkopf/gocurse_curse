package curse

import (
	"errors"
	"fmt"
)

type ErrNotFound struct {
	CurseID string
}

var ErrInvalidStartDate = errors.New("invalid start date")
var ErrInvalidEndDate = errors.New("invalid end date")
var ErrNameRequired = errors.New("name is required")
var ErrStartRequired = errors.New("start date is required")
var ErrEndRequired = errors.New("end date is required")
var ErrEndLesserStart = errors.New("end date mustn't be lesser than start date")

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("user '%s' not found", e.CurseID)
}
