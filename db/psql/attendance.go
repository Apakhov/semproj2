package psql

import (
	"errors"
	"semprojdb/queryer"
	"time"
)

type Attendance struct {
	ID        int64     // id
	Date      time.Time // date
	StudentID int64     // student_id
	TeacherID int64     // teacher_id
	CourseID  int64     // course_id
}

const sqlAttendanceInsert = `
INSERT INTO attendance (date, student_id, teacher_id, course_id)
VALUES ($1, $2, $3, $4) RETURNING id;`

func (db *DB) NewAttendance(v *Attendance) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlAttendanceInsert, v.Date, v.StudentID, v.TeacherID, v.CourseID).Scan(&v.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) UpdAttendance(v *Attendance) error {
	sql := queryer.Queryer{}
	sql.WriteString("UPDATE attendance SET")
	sql.WithSep(",")
	sql.IF(!v.Date.IsZero()).
		Sep().WriteStringArgs(" date = $$", v.Date)
	sql.IF(v.StudentID != -1).
		Sep().WriteStringArgs(" student_id = $$", v.StudentID)
	sql.IF(v.TeacherID != -1).
		Sep().WriteStringArgs(" teacher_id = $$", v.TeacherID)
	sql.IF(v.CourseID != -1).
		Sep().WriteStringArgs(" course_id = $$", v.CourseID)
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

func (db *DB) DelAttendance(id int64) error {
	sql := queryer.Queryer{}
	sql.WriteString("DELETE FROM attendance WHERE").WithSep(" AND ")
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

func (db *DB) GetAttendances(id int64, dateLE, dateGE time.Time, sID, tID, cID, limit int64) ([]Attendance, error) {
	limit = fixLimit(limit, 200)
	sql := queryer.Queryer{}
	sql.WriteStringArgs(`
		SELECT id, date, student_id, teacher_id, course_id
		FROM attendance WHERE id >= $$`, id)
	sql.
		IF(!dateLE.IsZero()).
		WriteStringArgs(" AND date <= $$", dateLE)
	sql.
		IF(!dateGE.IsZero()).
		WriteStringArgs(" AND date >= $$", dateGE)
	sql.
		IF(sID != -1).
		WriteStringArgs(" AND student_id = $$", sID)
	sql.
		IF(tID != -1).
		WriteStringArgs(" AND teacher_id = $$", tID)
	sql.
		IF(cID != -1).
		WriteStringArgs(" AND course_id = $$", cID)
	sql.
		WriteStringArgs(" ORDER BY id, date LIMIT $$", limit)

	ctx, cancel := defaultContext()
	defer cancel()
	db.lg.Info(sql.String(), sql.Args())
	rows, err := db.conn.Query(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return nil, err
	}
	vs := make([]Attendance, 0)
	for cnt := 0; rows.Next(); cnt++ {
		var v Attendance
		err := rows.Scan(&v.ID, &v.Date, &v.StudentID, &v.TeacherID, &v.CourseID)
		if err != nil {
			db.lg.Info("ERR2", err)
			return vs, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}
