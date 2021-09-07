package web

import (
	"fmt"
	"net/http"
)

func invalidOption(io string) error {
	return fmt.Errorf("%w: %s", ErrInvalidOption, io)
}

func writeError(w http.ResponseWriter, err int) {
	w.WriteHeader(err)
	w.Write([]byte(http.StatusText(err)))
}
