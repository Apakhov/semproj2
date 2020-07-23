package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
	"time"
)

type Course struct {
	ID        int64     `json:"id"`
	ShortName string    `json:"short_name"`
	FullName  string    `json:"full_name"`
	Semester  int64     `json:"semester"`
	BeginD    time.Time `json:"begin_d"`
	EndD      time.Time `json:"end_d"`
	SubjectID int64     `json:"subject_id"`
	StGroupID int64     `json:"st_group_id"`
	Active    bool      `json:"active"`
}

func courseToDB(s *Course, d *psql.Course) {
	d.ID = s.ID
	d.ShortName = s.ShortName
	d.FullName = s.FullName
	d.Semester = s.Semester
	d.BeginD = s.BeginD
	d.EndD = s.EndD
	d.SubjectID = s.SubjectID
	d.StGroupID = s.StGroupID
	d.Active = s.Active
}

func courseFromDB(s *psql.Course, d *Course) {
	d.ID = s.ID
	d.ShortName = s.ShortName
	d.FullName = s.FullName
	d.Semester = s.Semester
	d.BeginD = s.BeginD
	d.EndD = s.EndD
	d.SubjectID = s.SubjectID
	d.StGroupID = s.StGroupID
	d.Active = s.Active
}

func (h *Handler) HandleCourseCreate(w http.ResponseWriter, r *http.Request) {
	v := new(Course)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Course)
	courseToDB(v, vdb)
	err := h.db.NewCourse(vdb)
	courseFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleCourseUpdate(w http.ResponseWriter, r *http.Request) {
	v := new(Course)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Course)
	courseToDB(v, vdb)
	err := h.db.UpdCourse(vdb)
	courseFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleCourseRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	shName, _ := helpers.ReadGetString(r, "short_name")
	flName, _ := helpers.ReadGetString(r, "full_name")
	sem, _ := helpers.ReadGetInt64(r, "semester")
	begLE, _ := helpers.ReadGetTime(r, "beg_le")
	begGE, _ := helpers.ReadGetTime(r, "beg_ge")
	endLE, _ := helpers.ReadGetTime(r, "end_le")
	endGE, _ := helpers.ReadGetTime(r, "end_ge")
	subjID, _ := helpers.ReadGetInt64(r, "subject_id")
	stgrID, _ := helpers.ReadGetInt64(r, "st_group_id")
	active, _ := helpers.ReadGetString(r, "active")

	vdbs, err := h.db.GetCourses(id, shName, flName, sem, begLE, begGE, endLE, endGE, subjID, stgrID, active, limit)
	h.lg.Info(vdbs)
	vs := make([]Course, len(vdbs))
	for i := range vdbs {
		courseFromDB(&vdbs[i], &vs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vs,
	})
	w.Write(bt)
}

func (h *Handler) HandleCourseDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	err := h.db.DelCourse(id)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
