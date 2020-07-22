package psql

import (
	"errors"
	"semprojdb/queryer"
)

type Teacher struct {
	ID         int64
	ContractID string // contract_id
	FirsName   string // first_name
	LastName   string // last_name
	Email      string // email
	CathedraID int64  // cathedra_id
	Active     bool
}

const sqlTeachertInsert = `
INSERT INTO teacher (contract_id, first_name, last_name, email, cathedra_id) 
VALUES ($1, $2, $3, $4, $5) RETURNING id;`

func (db *DB) NewTeacher(v *Teacher) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlTeachertInsert, v.ContractID, v.FirsName, v.LastName, v.Email, v.CathedraID).Scan(&v.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) UpdTeacher(v *Teacher) error {
	sql := queryer.Queryer{}
	sql.WriteString("UPDATE teacher SET")
	sql.WithSep(",")
	sql.IF(v.ContractID != "").
		Sep().WriteStringArgs(" contract_id = $$", v.ContractID)
	sql.IF(v.FirsName != "").
		Sep().WriteStringArgs(" first_name = $$", v.FirsName)
	sql.IF(v.LastName != "").
		Sep().WriteStringArgs(" last_name = $$", v.LastName)
	sql.IF(v.Email != "").
		Sep().WriteStringArgs(" email = $$", v.Email)
	sql.IF(v.CathedraID != -1).
		Sep().WriteStringArgs(" cathedra_id = $$", v.CathedraID)
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

func (db *DB) DelTeacher(id int64, contractID, email string) error {
	sql := queryer.Queryer{}
	sql.WriteString("DELETE FROM teacher WHERE").WithSep(" AND ")
	sql.IF(contractID != "").
		Sep().WriteStringArgs(" contract_id = $$", contractID)
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

func (db *DB) GetTeachers(id int64, contID, fName, lName, email, active string, catID int64, limit int64) ([]Teacher, error) {
	limit = fixLimit(limit, 200)
	sql := queryer.Queryer{}
	sql.WriteStringArgs(`
		SELECT id, contract_id, first_name, last_name, email, cathedra_id, active
		FROM teacher WHERE id >= $$`, id)
	sql.
		IF(contID != "").
		WriteStringArgs(" AND contract_id = $$", contID)
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
		IF(catID != -1).
		WriteStringArgs(" AND cathedra_id = $$", catID)
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
