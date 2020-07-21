package psql

import (
	"context"
	"fmt"
	"semprojdb/logger"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type Config struct {
	Host     string
	Port     uint16
	Database string
	User     string
	Password string
}

type DB struct {
	lg         *logger.Logger
	connString string
	conn       *pgx.Conn
}

func NewDB(c *Config, lg *logger.Logger) *DB {
	if c.Port == 0 {
		c.Port = 5432
	}
	db := new(DB)
	db.lg = lg

	db.connString = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		c.Host,
		c.Port,
		c.Database,
		c.User,
		c.Password)
	return db
}

func (db *DB) Connect() error {
	var err error
	db.conn, err = pgx.Connect(context.Background(), db.connString)
	return errors.Wrap(err, "db connect")
}

var defaultTimeout = 2 * time.Second

func defaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}
