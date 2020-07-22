package psql

import (
	"errors"
	"semprojdb/queryer"
)

type Cathedra struct {
	ID        int64
	ShortName string
	FullName  string
	FacultyID int64
	IsDeleted bool
}

const sqlCathedraInsert = `
INSERT INTO cathedra (short_name, full_name, faculty_id) 
VALUES ($1, $2, $3) RETURNING id;`

func (db *DB) NewCathedra(c *Cathedra) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlCathedraInsert, c.ShortName, c.FullName, c.FacultyID).Scan(&c.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

func (db *DB) UpdCathedra(f *Cathedra) error {
	sql := queryer.Queryer{}
	sql.WriteString("UPDATE cathedra SET")
	sql.WithSep(",")
	sql.IF(f.ShortName != "").
		Sep().WriteStringArgs(" short_name = $$", f.ShortName)
	sql.IF(f.FullName != "").
		Sep().WriteStringArgs(" full_name = $$", f.FullName)
	sql.IF(f.FacultyID != -1).
		Sep().WriteStringArgs(" faculty_id = $$", f.FacultyID)
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

func (db *DB) DelCathedra(id int64, shName, flName string) error {
	sql := queryer.Queryer{}
	sql.WriteString("DELETE FROM cathedra WHERE").WithSep(" AND ")
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

func (db *DB) GetCathedras(id int64, shName, flName string, fid int64, limit int64) ([]Cathedra, error) {
	limit = fixLimit(limit, 50)
	sql := queryer.Queryer{}
	sql.WriteStringArgs(`
		SELECT id, short_name, full_name, faculty_id, deleted
		FROM cathedra WHERE id >= $$`, id)
	sql.
		IF(shName != "").
		WriteStringArgs(" AND short_name = $$", shName)
	sql.
		IF(flName != "").
		WriteStringArgs(" AND full_name = $$", flName)
	sql.
		IF(fid != -1).
		WriteStringArgs(" AND faculty_id = $$", fid)
	sql.WriteStringArgs(" ORDER BY id LIMIT $$", limit)

	ctx, cancel := defaultContext()
	defer cancel()
	db.lg.Info(sql.String(), sql.Args())
	rows, err := db.conn.Query(ctx, sql.String(), sql.Args()...)
	if err != nil {
		db.lg.Info("ERR1", err)
		return nil, err
	}
	fs := make([]Cathedra, 0)
	for cnt := 0; rows.Next(); cnt++ {
		var f Cathedra
		err := rows.Scan(&f.ID, &f.ShortName, &f.FullName, &f.FacultyID, &f.IsDeleted)
		if err != nil {
			db.lg.Info("ERR2", err)
			return fs, err
		}
		fs = append(fs, f)
	}
	return fs, nil
}
