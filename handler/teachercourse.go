package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
)

func (h *Handler) HandleAssignTeacherToCourse(w http.ResponseWriter, r *http.Request) {
	tID, _ := helpers.ReadGetInt64(r, "t_id")
	cID, _ := helpers.ReadGetInt64(r, "c_id")
	err := h.db.AssignTeacherToCourse(tID, cID)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}

func (h *Handler) HandleUnassignTeacherFromCourse(w http.ResponseWriter, r *http.Request) {
	tID, _ := helpers.ReadGetInt64(r, "t_id")
	cID, _ := helpers.ReadGetInt64(r, "c_id")
	err := h.db.UnassignTeacherFromCourse(tID, cID)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}

func (h *Handler) HandleGetTeachersFromCourse(w http.ResponseWriter, r *http.Request) {
	cID, _ := helpers.ReadGetInt64(r, "id")

	vdbs, err := h.db.GetTeachersFromCourse(cID)
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

func (h *Handler) HandleGetCoursesOfTeacher(w http.ResponseWriter, r *http.Request) {
	tID, _ := helpers.ReadGetInt64(r, "id")

	vdbs, err := h.db.GetCoursesOfTeacher(tID)
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
