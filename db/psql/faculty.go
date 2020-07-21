package psql

import (
	"errors"
	"semprojdb/queryer"
)

type Faculty struct {
	ID        int64
	ShortName string
	FullName  string
	IsDeleted bool
}

const sqlFacultyInsert = `
INSERT INTO faculty (short_name, full_name) 
VALUES ($1, $2) RETURNING id;`

func (db *DB) NewFaculty(fs *Faculty) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlFacultyInsert, fs.ShortName, fs.FullName).Scan(&fs.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) UpdFaculty(f *Faculty) error {
	sql := queryer.Queryer{}
	sql.WriteString("UPDATE faculty SET")
	sql.WithSep(",")
	sql.IF(f.ShortName != "").
		Sep().WriteStringArgs(" short_name = $$", f.ShortName)
	sql.IF(f.FullName != "").
		Sep().WriteStringArgs(" full_name = $$", f.FullName)
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

func (db *DB) DelFaculty(id int64, shName, flName string) error {
	sql := queryer.Queryer{}
	sql.WriteString("DELETE FROM faculty WHERE").WithSep(" AND ")
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

func (db *DB) GetFaculties(id int64, shName, flName string, limit int64) ([]Faculty, error) {
	limit = fixLimit(limit, 50)
	sql := queryer.Queryer{}
	sql.WriteStringArgs(`
		SELECT id, short_name, full_name, deleted
		FROM faculty WHERE id >= $$`, id)
	sql.
		IF(shName != "").
		WriteStringArgs(" AND short_name = $$", shName)
	sql.
		IF(flName != "").
		WriteStringArgs(" AND full_name = $$", flName)
	sql.WriteStringArgs(" ORDER BY id LIMIT $$", limit)

	ctx, cancel := defaultContext()
	defer cancel()
	db.lg.Info(sql.String(), sql.Args())
	rows, err := db.conn.Query(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return nil, err
	}
	fs := make([]Faculty, 0)
	for cnt := 0; rows.Next(); cnt++ {
		var f Faculty
		err := rows.Scan(&f.ID, &f.ShortName, &f.FullName, &f.IsDeleted)
		if err != nil {
			db.lg.Info("ERR2", err)
			return fs, err
		}
		fs = append(fs, f)
	}
	return fs, nil
}
