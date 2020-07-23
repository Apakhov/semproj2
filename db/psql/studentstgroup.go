package psql

import (
	"semprojdb/queryer"
)

const sqlStudentToGroup = `
INSERT INTO student_st_group (student_id, st_group_id) 
VALUES ($1, $2)`

func (db *DB) PutStudentToGroup(stID, grID int64) error {
	ctx, cancel := defaultContext()
	defer cancel()
	_, err := db.conn.Exec(ctx, sqlStudentToGroup, stID, grID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

const sqlGetStudentsFromGroup = `
SELECT s.id, s.study_id, s.first_name, s.last_name, s.email, s.active
FROM student s
         JOIN student_st_group sg on s.id = sg.student_id
WHERE sg.st_group_id = $$`

func (db *DB) GetStudentsFromGroup(grID int64, active string) ([]Student, error) {
	sql := queryer.Queryer{}
	sql.WriteStringArgs(sqlGetStudentsFromGroup, grID)
	sql.
		IF(active != "" && active == "t").
		WriteString(" AND sg.active = true")
	sql.
		IF(active != "" && active == "f").
		WriteString(" AND sg.active = false")

	ctx, cancel := defaultContext()
	defer cancel()
	db.lg.Info(sql.String(), sql.Args())
	rows, err := db.conn.Query(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return nil, err
	}
	vs := make([]Student, 0)
	for cnt := 0; rows.Next(); cnt++ {
		var v Student
		err := rows.Scan(&v.ID, &v.StudyID, &v.FirsName, &v.LastName, &v.Email, &v.Active)
		if err != nil {
			db.lg.Info("ERR2", err)
			return vs, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

const sqlGetGroupsOfStudent = `
SELECT gr.id, gr.group_id, gr.begin_d, gr.end_d, gr.teacher_id, gr.active
FROM st_group gr
         JOIN student_st_group sg on gr.id = sg.st_group_id
WHERE sg.student_id = $$`

func (db *DB) GetGroupsOfStudent(stID int64) ([]StGroup, error) {
	sql := queryer.Queryer{}
	sql.WriteStringArgs(sqlGetGroupsOfStudent, stID)

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
