package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
)

type Faculty struct {
	ID        int64  `json:"id"`
	ShortName string `json:"short_name"`
	FullName  string `json:"full_name"`
	IsDeleted bool   `json:"is_deleted"`
}

func facultyToDB(s *Faculty, d *psql.Faculty) {
	d.ID = s.ID
	d.ShortName = s.ShortName
	d.FullName = s.FullName
	d.IsDeleted = s.IsDeleted
}

func facultyFromDB(s *psql.Faculty, d *Faculty) {
	d.ID = s.ID
	d.ShortName = s.ShortName
	d.FullName = s.FullName
	d.IsDeleted = s.IsDeleted
}

func (h *Handler) HandleFacultyCreate(w http.ResponseWriter, r *http.Request) {
	f := new(Faculty)
	if !helpers.DecodeJSONBody(h.lg, w, r, f) {
		return
	}
	fdb := new(psql.Faculty)
	facultyToDB(f, fdb)
	err := h.db.NewFaculty(fdb)
	facultyFromDB(fdb, f)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: f,
	})
	w.Write(bt)
}

func (h *Handler) HandleFacultyUpdate(w http.ResponseWriter, r *http.Request) {
	f := new(Faculty)
	if !helpers.DecodeJSONBody(h.lg, w, r, f) {
		return
	}
	fdb := new(psql.Faculty)
	facultyToDB(f, fdb)
	err := h.db.UpdFaculty(fdb)
	facultyFromDB(fdb, f)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: f,
	})
	w.Write(bt)
}

func (h *Handler) HandleFacultyRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	sh, _ := helpers.ReadGetString(r, "short_name")
	fl, _ := helpers.ReadGetString(r, "full_name")
	h.lg.Info(id, sh, fl, limit)
	fdbs, err := h.db.GetFaculties(id, sh, fl, limit)
	h.lg.Info(fdbs)
	fs := make([]Faculty, len(fdbs))
	for i := range fdbs {
		facultyFromDB(&fdbs[i], &fs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: fs,
	})
	w.Write(bt)
}

func (h *Handler) HandleFacultyDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	sh, _ := helpers.ReadGetString(r, "short_name")
	fl, _ := helpers.ReadGetString(r, "full_name")
	err := h.db.DelFaculty(id, sh, fl)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
