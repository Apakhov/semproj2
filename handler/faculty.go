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
	err := h.db.NewFaculties(fdb)
	facultyFromDB(fdb, f)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: fdb,
	})
	w.Write(bt)
}
