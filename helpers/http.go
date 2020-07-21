package helpers

import (
	"errors"
	"net/http"
	"semprojdb/logger"
)

func HandleErr(lg *logger.Logger, w http.ResponseWriter, err error) bool {
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.msg, mr.status)
		} else {
			lg.Error(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return true
	}
	return false
}
