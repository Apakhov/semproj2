package psql

type Faculty struct {
	ID        int64
	ShortName string
	FullName  string
	IsDeleted bool
}

const sqlFacultyInsert = `
INSERT INTO faculty (short_name, full_name) 
VALUES ($1, $2) RETURNING id;`

func (db *DB) NewFaculties(fs *Faculty) error {
	ctx, cancel := defaultContext()
	defer cancel()
	err := db.conn.QueryRow(ctx, sqlFacultyInsert, fs.ShortName, fs.FullName).Scan(&fs.ID)
	if err != nil {
		db.lg.Info("ERR1", err)
		return err
	}
	return nil
}

// func (db *DB) NewFaculties(fs *Faculty) error {

// 	sql := strings.Builder{}
// 	sql.WriteString(sqlFacultyInsertBeg)
// 	args := make([]interface{}, 0, len(fs)*2)
// 	fslast := len(fs) - 1
// 	for i, f := range fs {
// 		sql.WriteString(sqlFacultyInsertVal)
// 		if i != fslast {
// 			sql.WriteByte(',')
// 		}
// 		args = append(args, f.ShortName, f.FullName)
// 	}
// 	// tx, err := db.conn.Begin(context.Background())
// 	// if err != nil {
// 	// 	return ErrInternal
// 	// }
// 	// defer tx.Commit(context.Background())
// 	db.lg.Info("REQUEST", args)
// 	ctx, cancel := defaultContext()
// 	defer cancel()
// 	rows, err := db.conn.Query(ctx, sql.String(), args...)
// 	if err != nil {
// 		db.lg.Info("ERR1", err)
// 		return err
// 	}

// 	for cnt := 0; rows.Next(); cnt++ {
// 		err := rows.Scan(&fs[cnt].ID)
// 		if err != nil {
// 			db.lg.Info("ERR2", err)
// 			return err
// 		}
// 	}
// 	db.lg.Info("KK", fs)
// 	return nil
// }
