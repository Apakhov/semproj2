package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
)

type Teacher struct {
	ID         int64  `json:"id"`
	ContractID string `json:"contract_id"`
	FirsName   string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	CathedraID int64  `json:"cathedra_id"`
	Active     bool   `json:"active"`
}

func teacherToDB(s *Teacher, d *psql.Teacher) {
	d.ID = s.ID
	d.ContractID = s.ContractID
	d.FirsName = s.FirsName
	d.LastName = s.LastName
	d.Email = s.Email
	d.CathedraID = s.CathedraID
	d.Active = s.Active
}

func teacherFromDB(s *psql.Teacher, d *Teacher) {
	d.ID = s.ID
	d.ContractID = s.ContractID
	d.FirsName = s.FirsName
	d.LastName = s.LastName
	d.Email = s.Email
	d.CathedraID = s.CathedraID
	d.Active = s.Active
}

func (h *Handler) HandleTeacherCreate(w http.ResponseWriter, r *http.Request) {
	v := new(Teacher)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Teacher)
	teacherToDB(v, vdb)
	err := h.db.NewTeacher(vdb)
	teacherFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleTeacherUpdate(w http.ResponseWriter, r *http.Request) {
	v := new(Teacher)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Teacher)
	teacherToDB(v, vdb)
	err := h.db.UpdTeacher(vdb)
	teacherFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleTeacherRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	contID, _ := helpers.ReadGetString(r, "contract_id")
	fName, _ := helpers.ReadGetString(r, "first_name")
	lName, _ := helpers.ReadGetString(r, "last_name")
	email, _ := helpers.ReadGetString(r, "email")
	cid, _ := helpers.ReadGetInt64(r, "cathedra_id")
	active, _ := helpers.ReadGetString(r, "active")
	vdbs, err := h.db.GetTeachers(id, contID, fName, lName, email, active, cid, limit)
	h.lg.Info(vdbs)
	vs := make([]Teacher, len(vdbs))
	for i := range vdbs {
		teacherFromDB(&vdbs[i], &vs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vs,
	})
	w.Write(bt)
}

func (h *Handler) HandleTeacherDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	cid, _ := helpers.ReadGetString(r, "contract_id")
	email, _ := helpers.ReadGetString(r, "email")
	err := h.db.DelTeacher(id, cid, email)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
