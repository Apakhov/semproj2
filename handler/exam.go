package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
	"time"
)

type Exam struct {
	ID        int64     `json:"id"`
	Date      time.Time `json:"date"`
	Points    int64     `json:"points"`
	Type      string    `json:"type"`
	StudentID int64     `json:"student_id"`
	TeacherID int64     `json:"teacher_id"`
	CourseID  int64     `json:"course_id"`
}

func examToDB(s *Exam, d *psql.Exam) {
	d.ID = s.ID
	d.Date = s.Date
	d.Points = s.Points
	d.Type = s.Type
	d.StudentID = s.StudentID
	d.TeacherID = s.TeacherID
	d.CourseID = s.CourseID
}

func examFromDB(s *psql.Exam, d *Exam) {
	d.ID = s.ID
	d.Date = s.Date
	d.Points = s.Points
	d.Type = s.Type
	d.StudentID = s.StudentID
	d.TeacherID = s.TeacherID
	d.CourseID = s.CourseID
}

func (h *Handler) HandleExamCreate(w http.ResponseWriter, r *http.Request) {
	v := new(Exam)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Exam)
	examToDB(v, vdb)
	err := h.db.NewExam(vdb)
	examFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleExamUpdate(w http.ResponseWriter, r *http.Request) {
	v := new(Exam)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Exam)
	examToDB(v, vdb)
	err := h.db.UpdExam(vdb)
	examFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleExamRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	dateLE, _ := helpers.ReadGetTime(r, "date_le")
	dateGE, _ := helpers.ReadGetTime(r, "date_ge")
	sID, _ := helpers.ReadGetInt64(r, "subject_id")
	tID, _ := helpers.ReadGetInt64(r, "teacher_id")
	cID, _ := helpers.ReadGetInt64(r, "course_id")
	tp, _ := helpers.ReadGetString(r, "type")

	vdbs, err := h.db.GetExams(id, dateLE, dateGE, sID, tID, cID, tp, limit)
	h.lg.Info(vdbs)
	vs := make([]Exam, len(vdbs))
	for i := range vdbs {
		examFromDB(&vdbs[i], &vs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vs,
	})
	w.Write(bt)
}

func (h *Handler) HandleExamDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	err := h.db.DelExam(id)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
