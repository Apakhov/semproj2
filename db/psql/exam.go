package psql

import (
	"errors"
	"semprojdb/queryer"
	"time"
)

type Exam struct {
	ID        int64     // id
	Date      time.Time // date
	Points    int64     // points
	Type      string    // type
	StudentID int64     // student_id
	TeacherID int64     // teacher_id
	CourseID  int64     // course_id
}

const sqlExamInsert = `
INSERT INTO exam (date, points, type, student_id, teacher_id, course_id)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

func (db *DB) NewExam(v *Exam) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlExamInsert, v.Date, v.Points, v.Type, v.StudentID, v.TeacherID, v.CourseID).Scan(&v.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) UpdExam(v *Exam) error {
	sql := queryer.Queryer{}
	sql.WriteString("UPDATE exam SET")
	sql.WithSep(",")
	sql.IF(!v.Date.IsZero()).
		Sep().WriteStringArgs(" date = $$", v.Date)
	sql.IF(v.Points != -1).
		Sep().WriteStringArgs(" points = $$", v.Points)
	sql.IF(v.Type != "").
		Sep().WriteStringArgs(" type = $$", v.Type)
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

func (db *DB) DelExam(id int64) error {
	sql := queryer.Queryer{}
	sql.WriteString("DELETE FROM exam WHERE").WithSep(" AND ")
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

func (db *DB) GetExams(id int64, dateLE, dateGE time.Time, sID, tID, cID int64, tp string, limit int64) ([]Exam, error) {
	limit = fixLimit(limit, 200)
	sql := queryer.Queryer{}
	sql.WriteStringArgs(`
		SELECT id, date, points, type, student_id, teacher_id, course_id
		FROM exam WHERE id >= $$`, id)
	sql.
		IF(!dateLE.IsZero()).
		WriteStringArgs(" AND date <= $$", dateLE)
	sql.
		IF(!dateGE.IsZero()).
		WriteStringArgs(" AND date >= $$", dateGE)
	sql.
		IF(tp != "").
		WriteStringArgs(" AND type = $$", tp)
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
	vs := make([]Exam, 0)
	for cnt := 0; rows.Next(); cnt++ {
		var v Exam
		err := rows.Scan(&v.ID, &v.Date, &v.Points, &v.Type, &v.StudentID, &v.TeacherID, &v.CourseID)
		if err != nil {
			db.lg.Info("ERR2", err)
			return vs, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}
