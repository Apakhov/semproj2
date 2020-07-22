package psql

import (
	"errors"
	"semprojdb/queryer"
	"time"
)

type StGroup struct {
	ID        int64     // id
	GroupID   string    // group_id
	BeginD    time.Time // begin_d
	EndD      time.Time // end_d
	TeacherID int64     // teacher_id
	Active    bool      // active
}

const sqlStGrouptInsert = `
INSERT INTO st_group (group_id, begin_d, end_d, teacher_id)
VALUES ($1, $2, $3, $4) RETURNING id;`

func (db *DB) NewStGroup(v *StGroup) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlStGrouptInsert, v.GroupID, v.BeginD, v.EndD, v.TeacherID).Scan(&v.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) UpdStGroup(v *StGroup) error {
	sql := queryer.Queryer{}
	sql.WriteString("UPDATE st_group SET")
	sql.WithSep(",")
	sql.IF(v.GroupID != "").
		Sep().WriteStringArgs(" group_id = $$", v.GroupID)
	sql.IF(!v.BeginD.IsZero()).
		Sep().WriteStringArgs(" begin_d = $$", v.BeginD)
	sql.IF(!v.EndD.IsZero()).
		Sep().WriteStringArgs(" end_d = $$", v.EndD)
	sql.IF(v.TeacherID != -1).
		Sep().WriteStringArgs(" teacher_id = $$", v.TeacherID)
	sql.WriteStringArgs(" WHERE id = $$", v.ID)

	ctx, cancel := defaultContext()
	defer cancel()
	_, err := db.conn.Exec(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) DelStGroup(id int64, groupID string) error {
	sql := queryer.Queryer{}
	sql.WriteString("DELETE FROM st_group WHERE").WithSep(" AND ")
	sql.IF(groupID != "").
		Sep().WriteStringArgs(" group_id = $$", groupID)
	sql.IF(id != -1).
		Sep().WriteStringArgs(" id = $$", id)
	if len(sql.Args()) == 0 {
		return errors.New("cant delete nothing")
	}

	ctx, cancel := defaultContext()
	defer cancel()
	_, err := db.conn.Exec(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) GetStGroups(
	id int64, grID string, tID int64,
	begLE, begGE, endLE, endGE time.Time,
	limit int64,
) ([]StGroup, error) {
	limit = fixLimit(limit, 200)
	sql := queryer.Queryer{}
	sql.WriteStringArgs(`
		SELECT id, group_id, begin_d, end_d, teacher_id, active
		FROM st_group WHERE id >= $$`, id)
	sql.
		IF(grID != "").
		WriteStringArgs(" AND group_id = $$", grID)
	sql.
		IF(tID != -1).
		WriteStringArgs(" AND teacher_id = $$", tID)
	sql.
		IF(!begLE.IsZero()).
		WriteStringArgs(" AND begin_d <= $$", begLE)
	sql.
		IF(!begGE.IsZero()).
		WriteStringArgs(" AND begin_d >= $$", begGE)
	sql.
		IF(!endLE.IsZero()).
		WriteStringArgs(" AND end_d <= $$", endLE)
	sql.
		IF(!endGE.IsZero()).
		WriteStringArgs(" AND end_d >= $$", endGE)
	sql.
		WriteStringArgs(" ORDER BY id LIMIT $$", limit)

	ctx, cancel := defaultContext()
	defer cancel()
	db.lg.Info(sql.String(), sql.Args())
	rows, err := db.conn.Query(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return nil, err
	}
	vs := make([]StGroup, 0)
	for cnt := 0; rows.Next(); cnt++ {
		var v StGroup
		err := rows.Scan(&v.ID, &v.GroupID, &v.BeginD, &v.EndD, &v.TeacherID, &v.Active)
		if err != nil {
			db.lg.Info("ERR2", err)
			return vs, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}
