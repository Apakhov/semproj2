package helpers

import (
	"errors"
	"net/http"
	"semprojdb/logger"
	"strconv"
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

func ReadGetInt64(r *http.Request, name string) (int64, bool) {
	keys, ok := r.URL.Query()[name]
	if !ok || len(keys[0]) < 1 {
		return -1, false
	}
	res, err := strconv.ParseInt(keys[0], 0, 64)
	if err != nil {
		return -1, false
	}
	return res, true
}

func ReadGetString(r *http.Request, name string) (string, bool) {
	keys, ok := r.URL.Query()[name]
	if !ok || len(keys[0]) < 1 {
		return "", false
	}
	return keys[0], true
}
