package psql

import "github.com/pkg/errors"

var ErrInternal = errors.New("internal db error")

func AnalizeDBError(err error) string {
	if err == nil {
		return ""
	}
	return errors.Wrap(err, "dberror").Error()
}
