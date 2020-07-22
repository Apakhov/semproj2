package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
)

func (h *Handler) HandlePutStudentToGroup(w http.ResponseWriter, r *http.Request) {
	grID, _ := helpers.ReadGetInt64(r, "gr_id")
	stID, _ := helpers.ReadGetInt64(r, "st_id")
	err := h.db.PutStudentToGroup(stID, grID)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}

func (h *Handler) HandleGetStudentsFromGroup(w http.ResponseWriter, r *http.Request) {
	grID, _ := helpers.ReadGetInt64(r, "id")
	active, _ := helpers.ReadGetString(r, "active")

	vdbs, err := h.db.GetStudentsFromGroup(grID, active)
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

func (h *Handler) HandleGetGroupsOfStudent(w http.ResponseWriter, r *http.Request) {
	stID, _ := helpers.ReadGetInt64(r, "id")

	vdbs, err := h.db.GetGroupsOfStudent(stID)
	vs := make([]StGroup, len(vdbs))
	for i := range vdbs {
		stGroupFromDB(&vdbs[i], &vs[i])
	}
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vs,
	})
	w.Write(bt)
}
