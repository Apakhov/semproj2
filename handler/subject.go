package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
)

type Subject struct {
	ID        int64  `json:"id"`
	ShortName string `json:"short_name"`
	FullName  string `json:"full_name"`
	CathedraID int64  `json:"cathedra_id"`
	IsDeleted bool   `json:"is_deleted"`
}

func subjectToDB(s *Subject, d *psql.Subject) {
	d.ID = s.ID
	d.ShortName = s.ShortName
	d.FullName = s.FullName
	d.CathedraID = s.CathedraID
	d.IsDeleted = s.IsDeleted
}

func subjectFromDB(s *psql.Subject, d *Subject) {
	d.ID = s.ID
	d.ShortName = s.ShortName
	d.FullName = s.FullName
	d.CathedraID = s.CathedraID
	d.IsDeleted = s.IsDeleted
}

func (h *Handler) HandleSubjectCreate(w http.ResponseWriter, r *http.Request) {
	f := new(Subject)
	if !helpers.DecodeJSONBody(h.lg, w, r, f) {
		return
	}
	fdb := new(psql.Subject)
	subjectToDB(f, fdb)
	err := h.db.NewSubject(fdb)
	subjectFromDB(fdb, f)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: fdb,
	})
	w.Write(bt)
}

func (h *Handler) HandleSubjectUpdate(w http.ResponseWriter, r *http.Request) {
	f := new(Subject)
	if !helpers.DecodeJSONBody(h.lg, w, r, f) {
		return
	}
	fdb := new(psql.Subject)
	subjectToDB(f, fdb)
	err := h.db.UpdSubject(fdb)
	subjectFromDB(fdb, f)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: fdb,
	})
	w.Write(bt)
}

func (h *Handler) HandleSubjectRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	sh, _ := helpers.ReadGetString(r, "short_name")
	fl, _ := helpers.ReadGetString(r, "full_name")
	fid, _ := helpers.ReadGetInt64(r, "cathedra_id")
	h.lg.Info(id, sh, fl, limit)
	fdbs, err := h.db.GetSubjects(id, sh, fl, fid, limit)
	h.lg.Info(fdbs)
	fs := make([]Subject, len(fdbs))
	for i := range fdbs {
		subjectFromDB(&fdbs[i], &fs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: fs,
	})
	w.Write(bt)
}

func (h *Handler) HandleSubjectDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	sh, _ := helpers.ReadGetString(r, "short_name")
	fl, _ := helpers.ReadGetString(r, "full_name")
	err := h.db.DelSubject(id, sh, fl)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
