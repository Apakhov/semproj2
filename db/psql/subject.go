package psql

import (
	"errors"
	"semprojdb/queryer"
)

type Subject struct {
	ID         int64
	ShortName  string
	FullName   string
	CathedraID int64
	IsDeleted  bool
}

const sqlSubjectInsert = `
INSERT INTO subject (short_name, full_name, cathedra_id) 
VALUES ($1, $2, $3) RETURNING id;`

func (db *DB) NewSubject(c *Subject) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlSubjectInsert, c.ShortName, c.FullName, c.CathedraID).Scan(&c.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) UpdSubject(f *Subject) error {
	sql := queryer.Queryer{}
	sql.WriteString("UPDATE subject SET")
	sql.WithSep(",")
	sql.IF(f.ShortName != "").
		Sep().WriteStringArgs(" short_name = $$", f.ShortName)
	sql.IF(f.FullName != "").
		Sep().WriteStringArgs(" full_name = $$", f.FullName)
	sql.IF(f.CathedraID != -1).
		Sep().WriteStringArgs(" cathedra_id = $$", f.CathedraID)
	sql.WriteStringArgs(" WHERE id = $$", f.ID)

	ctx, cancel := defaultContext()
	defer cancel()
	_, err := db.conn.Exec(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) DelSubject(id int64, shName, flName string) error {
	sql := queryer.Queryer{}
	sql.WriteString("DELETE FROM subject WHERE").WithSep(" AND ")
	sql.IF(shName != "").
		Sep().WriteStringArgs(" short_name = $$", shName)
	sql.IF(flName != "").
		Sep().WriteStringArgs(" full_name = $$", flName)
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

func (db *DB) GetSubjects(id int64, shName, flName string, fid int64, limit int64) ([]Subject, error) {
	limit = fixLimit(limit, 50)
	sql := queryer.Queryer{}
	sql.WriteStringArgs(`
		SELECT id, short_name, full_name, cathedra_id, deleted
		FROM subject WHERE id >= $$`, id)
	sql.
		IF(shName != "").
		WriteStringArgs(" AND short_name = $$", shName)
	sql.
		IF(flName != "").
		WriteStringArgs(" AND full_name = $$", flName)
	sql.
		IF(fid != -1).
		WriteStringArgs(" AND cathedra_id = $$", fid)
	sql.WriteStringArgs(" ORDER BY id LIMIT $$", limit)

	ctx, cancel := defaultContext()
	defer cancel()
	db.lg.Info(sql.String(), sql.Args())
	rows, err := db.conn.Query(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return nil, err
	}
	fs := make([]Subject, 0)
	for cnt := 0; rows.Next(); cnt++ {
		var f Subject
		err := rows.Scan(&f.ID, &f.ShortName, &f.FullName, &f.CathedraID, &f.IsDeleted)
		if err != nil {
			db.lg.Info("ERR2", err)
			return fs, err
		}
		fs = append(fs, f)
	}
	return fs, nil
}
