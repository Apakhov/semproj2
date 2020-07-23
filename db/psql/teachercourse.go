package psql

import (
	"semprojdb/queryer"
)

const sqlATeacherToCourse = `
INSERT INTO teacher_course (teacher_id, course_id) 
VALUES ($1, $2)`

func (db *DB) AssignTeacherToCourse(tID, cID int64) error {
	ctx, cancel := defaultContext()
	defer cancel()
	_, err := db.conn.Exec(ctx, sqlATeacherToCourse, tID, cID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

const sqlUTeacherToCourse = `
DELETE FROM teacher_course 
WHERE (teacher_id, course_id) = ($1, $2)`

func (db *DB) UnassignTeacherFromCourse(tID, cID int64) error {
	ctx, cancel := defaultContext()
	defer cancel()
	_, err := db.conn.Exec(ctx, sqlUTeacherToCourse, tID, cID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

const sqlGetTeachersFromCourse = `
SELECT t.id, t.contract_id, t.first_name, t.last_name, t.email, t.cathedra_id, t.active
FROM teacher t
         JOIN teacher_course tc on t.id = tc.teacher_id
WHERE tc.course_id = $$`

func (db *DB) GetTeachersFromCourse(courID int64) ([]Teacher, error) {
	sql := queryer.Queryer{}
	sql.WriteStringArgs(sqlGetTeachersFromCourse, courID)

	ctx, cancel := defaultContext()
	defer cancel()
	db.lg.Info(sql.String(), sql.Args())
	rows, err := db.conn.Query(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return nil, err
	}
	vs := make([]Teacher, 0)
	for cnt := 0; rows.Next(); cnt++ {
		var v Teacher
		err := rows.Scan(&v.ID, &v.ContractID, &v.FirsName, &v.LastName, &v.Email, &v.CathedraID, &v.Active)
		if err != nil {
			db.lg.Info("ERR2", err)
			return vs, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

const sqlGetCoursesOfTeacher = `
SELECT c.id, c.short_name, c.full_name, c.semester, c.begin_d, c.end_d, c.subject_id, c.st_group_id, c.active
FROM course c
         JOIN teacher_course tc on c.id = tc.course_id
WHERE tc.teacher_id = $$`

func (db *DB) GetCoursesOfTeacher(tID int64) ([]Course, error) {
	sql := queryer.Queryer{}
	sql.WriteStringArgs(sqlGetCoursesOfTeacher, tID)

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
