package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
)

type Cathedra struct {
	ID        int64  `json:"id"`
	ShortName string `json:"short_name"`
	FullName  string `json:"full_name"`
	FacultyID int64  `json:"faculty_id"`
	IsDeleted bool   `json:"is_deleted"`
}

func cathedraToDB(s *Cathedra, d *psql.Cathedra) {
	d.ID = s.ID
	d.ShortName = s.ShortName
	d.FullName = s.FullName
	d.FacultyID = s.FacultyID
	d.IsDeleted = s.IsDeleted
}

func cathedraFromDB(s *psql.Cathedra, d *Cathedra) {
	d.ID = s.ID
	d.ShortName = s.ShortName
	d.FullName = s.FullName
	d.FacultyID = s.FacultyID
	d.IsDeleted = s.IsDeleted
}

func (h *Handler) HandleCathedraCreate(w http.ResponseWriter, r *http.Request) {
	f := new(Cathedra)
	if !helpers.DecodeJSONBody(h.lg, w, r, f) {
		return
	}
	fdb := new(psql.Cathedra)
	cathedraToDB(f, fdb)
	err := h.db.NewCathedra(fdb)
	cathedraFromDB(fdb, f)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: f,
	})
	w.Write(bt)
}

func (h *Handler) HandleCathedraUpdate(w http.ResponseWriter, r *http.Request) {
	f := new(Cathedra)
	if !helpers.DecodeJSONBody(h.lg, w, r, f) {
		return
	}
	fdb := new(psql.Cathedra)
	cathedraToDB(f, fdb)
	err := h.db.UpdCathedra(fdb)
	cathedraFromDB(fdb, f)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: f,
	})
	w.Write(bt)
}

func (h *Handler) HandleCathedraRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	sh, _ := helpers.ReadGetString(r, "short_name")
	fl, _ := helpers.ReadGetString(r, "full_name")
	fid, _ := helpers.ReadGetInt64(r, "faculty_id")
	h.lg.Info(id, sh, fl, limit)
	fdbs, err := h.db.GetCathedras(id, sh, fl, fid, limit)
	h.lg.Info(fdbs)
	fs := make([]Cathedra, len(fdbs))
	for i := range fdbs {
		cathedraFromDB(&fdbs[i], &fs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: fs,
	})
	w.Write(bt)
}

func (h *Handler) HandleCathedraDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	sh, _ := helpers.ReadGetString(r, "short_name")
	fl, _ := helpers.ReadGetString(r, "full_name")
	err := h.db.DelCathedra(id, sh, fl)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
