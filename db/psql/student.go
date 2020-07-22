package psql

import (
	"errors"
	"semprojdb/queryer"
)

type Student struct {
	ID       int64  // id
	StudyID  string // study_id
	FirsName string // first_name
	LastName string // last_name
	Email    string // email
	Active   bool   // active
}

const sqlStudenttInsert = `
INSERT INTO student (study_id, first_name, last_name, email) 
VALUES ($1, $2, $3, $4) RETURNING id;`

func (db *DB) NewStudent(v *Student) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlStudenttInsert, v.StudyID, v.FirsName, v.LastName, v.Email).Scan(&v.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) UpdStudent(v *Student) error {
	sql := queryer.Queryer{}
	sql.WriteString("UPDATE student SET")
	sql.WithSep(",")
	sql.IF(v.StudyID != "").
		Sep().WriteStringArgs(" study_id = $$", v.StudyID)
	sql.IF(v.FirsName != "").
		Sep().WriteStringArgs(" first_name = $$", v.FirsName)
	sql.IF(v.LastName != "").
		Sep().WriteStringArgs(" last_name = $$", v.LastName)
	sql.IF(v.Email != "").
		Sep().WriteStringArgs(" email = $$", v.Email)
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

func (db *DB) DelStudent(id int64, contractID, email string) error {
	sql := queryer.Queryer{}
	sql.WriteString("DELETE FROM student WHERE").WithSep(" AND ")
	sql.IF(contractID != "").
		Sep().WriteStringArgs(" study_id = $$", contractID)
	sql.IF(email != "").
		Sep().WriteStringArgs(" email = $$", email)
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

func (db *DB) GetStudents(id int64, stID, fName, lName, email, active string, limit int64) ([]Student, error) {
	limit = fixLimit(limit, 200)
	sql := queryer.Queryer{}
	sql.WriteStringArgs(`
		SELECT id, study_id, first_name, last_name, email, active
		FROM student WHERE id >= $$`, id)
	sql.
		IF(stID != "").
		WriteStringArgs(" AND study_id = $$", stID)
	sql.
		IF(fName != "").
		WriteStringArgs(" AND first_name = $$", fName)
	sql.
		IF(lName != "").
		WriteStringArgs(" AND last_name = $$", lName)
	sql.
		IF(email != "").
		WriteStringArgs(" AND email = $$", email)
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
