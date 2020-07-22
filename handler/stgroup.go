package handler

import (
	"encoding/json"
	"net/http"
	"semprojdb/db/psql"
	"semprojdb/helpers"
	"time"
)

type StGroup struct {
	ID        int64     `json:"id"`
	GroupID   string    `json:"group_id"`
	BeginD    time.Time `json:"begin_d"`
	EndD      time.Time `json:"end_d"`
	TeacherID int64     `json:"teacher_id"`
	Active    bool      `json:"active"`
}

func stGroupToDB(s *StGroup, d *psql.StGroup) {
	d.ID = s.ID
	d.GroupID = s.GroupID
	d.BeginD = s.BeginD
	d.EndD = s.EndD
	d.TeacherID = s.TeacherID
	d.Active = s.Active
}

func stGroupFromDB(s *psql.StGroup, d *StGroup) {
	d.ID = s.ID
	d.GroupID = s.GroupID
	d.BeginD = s.BeginD
	d.EndD = s.EndD
	d.TeacherID = s.TeacherID
	d.Active = s.Active
}

func (h *Handler) HandleStGroupCreate(w http.ResponseWriter, r *http.Request) {
	v := new(StGroup)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.StGroup)
	stGroupToDB(v, vdb)
	err := h.db.NewStGroup(vdb)
	stGroupFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vdb,
	})
	w.Write(bt)
}

func (h *Handler) HandleStGroupUpdate(w http.ResponseWriter, r *http.Request) {
	v := new(StGroup)
	if !helpers.DecodeJSONBody(h.lg, w, r, v) {
		return
	}
	vdb := new(psql.StGroup)
	stGroupToDB(v, vdb)
	err := h.db.UpdStGroup(vdb)
	stGroupFromDB(vdb, v)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: vdb,
	})
	w.Write(bt)
}

func (h *Handler) HandleStGroupRead(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	limit, lok := helpers.ReadGetInt64(r, "limit")
	if !lok {
		limit = 1
	}
	grID, _ := helpers.ReadGetString(r, "group_id")
	tID, _ := helpers.ReadGetInt64(r, "teacher_id")
	begLE, _ := helpers.ReadGetTime(r, "beg_le")
	begGE, _ := helpers.ReadGetTime(r, "beg_ge")
	endLE, _ := helpers.ReadGetTime(r, "end_le")
	endGE, _ := helpers.ReadGetTime(r, "end_ge")

	vdbs, err := h.db.GetStGroups(id, grID, tID, begLE, begGE, endLE, endGE, limit)
	h.lg.Info(vdbs)
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

func (h *Handler) HandleStGroupDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := helpers.ReadGetInt64(r, "id")
	grID, _ := helpers.ReadGetString(r, "group_id")
	err := h.db.DelStGroup(id, grID)
	bt, _ := json.Marshal(WithError{
		Err: psql.AnalizeDBError(err),
		Val: nil,
	})
	w.Write(bt)
}
