package psql

import (
	"errors"
	"semprojdb/queryer"
	"time"
)

type Course struct {
	ID        int64     // id
	ShortName string    // short_name
	FullName  string    // full_name
	Semester  int64     // semester
	BeginD    time.Time // begin_d
	EndD      time.Time // end_d
	SubjectID int64     // subject_id
	StGroupID int64     // st_group_id
	Active    bool      // active
}

// id, short_name, full_name, semester, begin_d, end_d, subject_id, st_group_id, active
const sqlCourseInsert = `
INSERT INTO course (short_name, full_name, semester, begin_d, end_d, subject_id, st_group_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

func (db *DB) NewCourse(v *Course) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlCourseInsert,
		v.ShortName, v.FullName, v.Semester, v.BeginD, v.EndD, v.SubjectID, v.StGroupID).Scan(&v.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	v.Active = true
	return nil
}

func (db *DB) UpdCourse(v *Course) error {
	sql := queryer.Queryer{}
	sql.WriteString("UPDATE course SET")
	sql.WithSep(",")
	sql.IF(v.ShortName != "").
		Sep().WriteStringArgs(" short_name = $$", v.ShortName)
	sql.IF(v.FullName != "").
		Sep().WriteStringArgs(" full_name = $$", v.FullName)
	sql.IF(v.Semester != -1).
		Sep().WriteStringArgs(" semester = $$", v.Semester)
	sql.IF(!v.BeginD.IsZero()).
		Sep().WriteStringArgs(" begin_d = $$", v.BeginD)
	sql.IF(!v.EndD.IsZero()).
		Sep().WriteStringArgs(" end_d = $$", v.EndD)
	sql.IF(v.SubjectID != -1).
		Sep().WriteStringArgs(" subject_id = $$", v.SubjectID)
	sql.IF(v.StGroupID != -1).
		Sep().WriteStringArgs(" st_group_id = $$", v.StGroupID)
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

func (db *DB) DelCourse(id int64) error {
	sql := queryer.Queryer{}
	sql.WriteString("DELETE FROM course WHERE").WithSep(" AND ")
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

func (db *DB) GetCourses(
	id int64, shName, flName string, sem int64,
	begLE, begGE, endLE, endGE time.Time,
	subjID, stgrID int64, active string,
	limit int64,
) ([]Course, error) {
	limit = fixLimit(limit, 200)
	sql := queryer.Queryer{}
	sql.WriteStringArgs(`
		SELECT id, short_name, full_name, semester, begin_d, end_d, subject_id, st_group_id, active
		FROM course WHERE id >= $$`, id)
	sql.
		IF(shName != "").
		WriteStringArgs(" AND short_name = $$", shName)
	sql.
		IF(flName != "").
		WriteStringArgs(" AND full_name = $$", flName)
	sql.
		IF(sem != -1).
		WriteStringArgs(" AND semester = $$", sem)
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
		IF(subjID != -1).
		WriteStringArgs(" AND subject_id = $$", subjID)
	sql.
		IF(stgrID != -1).
		WriteStringArgs(" AND st_group_id = $$", stgrID)
	sql.
		IF(active != "" && active == "t").
		WriteString(" AND active = true")
	sql.
		IF(active != "" && active == "f").
		WriteString(" AND active = false")
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
	vs := make([]Course, 0)
	for cnt := 0; rows.Next(); cnt++ {
		var v Course
		err := rows.Scan(&v.ID, &v.ShortName, &v.FullName, &v.Semester,
			&v.BeginD, &v.EndD, &v.SubjectID, &v.StGroupID, &v.Active)
		if err != nil {
			db.lg.Info("ERR2", err)
			return vs, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}
