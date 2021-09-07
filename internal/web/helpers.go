package web

import (
	"fmt"
)

func invalidOption(io string) error {
	return fmt.Errorf("%w: %s", ErrInvalidOption, io)
}
