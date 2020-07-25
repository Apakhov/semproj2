package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
	"time"
)

type Mark struct {
	ID        int64     `json:"id"`
	Date      time.Time `json:"date"`
	Points    int64     `json:"points"`
	StudentID int64     `json:"student_id"`
	TeacherID int64     `json:"teacher_id"`
	CourseID  int64     `json:"course_id"`
}

func markToDB(s *Mark, d *psql.Mark) {
	d.ID = s.ID
	d.Date = s.Date
	d.Points = s.Points
	d.StudentID = s.StudentID
	d.TeacherID = s.TeacherID
	d.CourseID = s.CourseID
}

func markFromDB(s *psql.Mark, d *Mark) {
	d.ID = s.ID
	d.Date = s.Date
	d.Points = s.Points
	d.StudentID = s.StudentID
	d.TeacherID = s.TeacherID
	d.CourseID = s.CourseID
}

func (h *Handler) HandleMarkCreate(w http.ResponseWriter, r *http.Request) {
	v := new(Mark)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Mark)
	markToDB(v, vdb)
	err := h.db.NewMark(vdb)
	markFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleMarkUpdate(w http.ResponseWriter, r *http.Request) {
	v := new(Mark)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Mark)
	markToDB(v, vdb)
	err := h.db.UpdMark(vdb)
	markFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleMarkRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	dateLE, _ := helpers.ReadGetTime(r, "date_le")
	dateGE, _ := helpers.ReadGetTime(r, "date_ge")
	sID, _ := helpers.ReadGetInt64(r, "student_id")
	tID, _ := helpers.ReadGetInt64(r, "teacher_id")
	cID, _ := helpers.ReadGetInt64(r, "course_id")

	vdbs, err := h.db.GetMarks(id, dateLE, dateGE, sID, tID, cID, limit)
	h.lg.Info(vdbs)
	vs := make([]Mark, len(vdbs))
	for i := range vdbs {
		markFromDB(&vdbs[i], &vs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vs,
	})
	w.Write(bt)
}

func (h *Handler) HandleMarkDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	err := h.db.DelMark(id)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
