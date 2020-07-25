package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
	"time"
)

type Attendance struct {
	ID        int64     `json:"id"`
	Date      time.Time `json:"date"`
	StudentID int64     `json:"student_id"`
	TeacherID int64     `json:"teacher_id"`
	CourseID  int64     `json:"course_id"`
}

func attendanceToDB(s *Attendance, d *psql.Attendance) {
	d.ID = s.ID
	d.Date = s.Date
	d.StudentID = s.StudentID
	d.TeacherID = s.TeacherID
	d.CourseID = s.CourseID
}

func attendanceFromDB(s *psql.Attendance, d *Attendance) {
	d.ID = s.ID
	d.Date = s.Date
	d.StudentID = s.StudentID
	d.TeacherID = s.TeacherID
	d.CourseID = s.CourseID
}

func (h *Handler) HandleAttendanceCreate(w http.ResponseWriter, r *http.Request) {
	v := new(Attendance)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Attendance)
	attendanceToDB(v, vdb)
	err := h.db.NewAttendance(vdb)
	attendanceFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleAttendanceUpdate(w http.ResponseWriter, r *http.Request) {
	v := new(Attendance)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.Attendance)
	attendanceToDB(v, vdb)
	err := h.db.UpdAttendance(vdb)
	attendanceFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: v,
	})
	w.Write(bt)
}

func (h *Handler) HandleAttendanceRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	dateLE, _ := helpers.ReadGetTime(r, "beg_le")
	dateGE, _ := helpers.ReadGetTime(r, "beg_ge")
	sID, _ := helpers.ReadGetInt64(r, "student_id")
	tID, _ := helpers.ReadGetInt64(r, "teacher_id")
	cID, _ := helpers.ReadGetInt64(r, "course_id")

	vdbs, err := h.db.GetAttendances(id, dateLE, dateGE, sID, tID, cID, limit)
	h.lg.Info(vdbs)
	vs := make([]Attendance, len(vdbs))
	for i := range vdbs {
		attendanceFromDB(&vdbs[i], &vs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vs,
	})
	w.Write(bt)
}

func (h *Handler) HandleAttendanceDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	err := h.db.DelAttendance(id)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
