package queryer

import (
	"fmt"
	"strconv"
	"strings"
)

type Queryer struct {
	b        strings.Builder
	args     []interface{}
	sep      string
	sepInUse bool
}

func (q *Queryer) String() string {
	return q.b.String()
}
func (q *Queryer) Args() []interface{} {
	return q.args
}
func (q *Queryer) Write(p []byte) *Queryer {
	if q == nil {
		return nil
	}
	q.b.Write(p)
	return q
}
func (q *Queryer) WriteString(s string) *Queryer {
	if q == nil {
		return nil
	}
	q.b.WriteString(s)
	return q
}
func (q *Queryer) WriteRune(r rune) *Queryer {
	if q == nil {
		return nil
	}
	q.b.WriteRune(r)
	return q
}
func (q *Queryer) WriteByte(b byte) *Queryer {
	if q == nil {
		return nil
	}
	q.b.WriteByte(b)
	return q
}
func (q *Queryer) Reset() {
	q.b.Reset()
}

func (q *Queryer) WriteStringArgs(s string, args ...interface{}) *Queryer {
	if q == nil {
		return nil
	}
	beg := 0
	ind := 0
	cnt := 0
	for ind != -1 {
		ind = strings.Index(s[beg:], "$$")
		if ind == -1 {
			break
		}
		q.WriteString(s[beg : beg+ind])
		q.WriteByte('$')
		q.WriteString(strconv.FormatInt(int64(len(q.args)+cnt)+1, 10))
		beg = beg + ind + 2
		cnt++
	}
	if cnt != len(args) {
		panic(fmt.Sprint("cant format ", s, args))
	}
	q.args = append(q.args, args...)
	return q
}

func (q *Queryer) WithSep(s string) *Queryer {
	if q == nil {
		return nil
	}
	q.sep = s
	q.sepInUse = false
	return q
}

func (q *Queryer) Sep() *Queryer {
	if q == nil {
		return nil
	}
	if !q.sepInUse {
		q.sepInUse = true
		return q
	}
	q.WriteString(q.sep)
	return q
}

func (q *Queryer) IF(b bool) *Queryer {
	if b {
		return q
	}
	return nil
}
