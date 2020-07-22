package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
)

type Student struct {
	ID       int64  `json:"id"`
	StudyID  string `json:"study_id"`
	FirsName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}

func studentToDB(s *Student, d *psql.Student) {
	d.ID = s.ID
	d.StudyID = s.StudyID
	d.FirsName = s.FirsName
	d.LastName = s.LastName
	d.Email = s.Email
	d.Active = s.Active
}

func studentFromDB(s *psql.Student, d *Student) {
	d.ID = s.ID
	d.StudyID = s.StudyID
	d.FirsName = s.FirsName
	d.LastName = s.LastName
	d.Email = s.Email
	d.Active = s.Active
}

func (h *Handler) HandleStudentCreate(w http.ResponseWriter, r *http.Request) {
	v := new(Student)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Student)
	studentToDB(v, vdb)
	err := h.db.NewStudent(vdb)
	studentFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vdb,
	})
	w.Write(bt)
}

func (h *Handler) HandleStudentUpdate(w http.ResponseWriter, r *http.Request) {
	v := new(Student)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Student)
	studentToDB(v, vdb)
	err := h.db.UpdStudent(vdb)
	studentFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vdb,
	})
	w.Write(bt)
}

func (h *Handler) HandleStudentRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	stID, _ := helpers.ReadGetString(r, "study_id")
	fName, _ := helpers.ReadGetString(r, "first_name")
	lName, _ := helpers.ReadGetString(r, "last_name")
	email, _ := helpers.ReadGetString(r, "email")
	active, _ := helpers.ReadGetString(r, "active")
	vdbs, err := h.db.GetStudents(id, stID, fName, lName, email, active, limit)
	h.lg.Info(vdbs)
	vs := make([]Student, len(vdbs))
	for i := range vdbs {
		studentFromDB(&vdbs[i], &vs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vs,
	})
	w.Write(bt)
}

func (h *Handler) HandleStudentDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	cid, _ := helpers.ReadGetString(r, "study_id")
	email, _ := helpers.ReadGetString(r, "email")
	err := h.db.DelStudent(id, cid, email)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
