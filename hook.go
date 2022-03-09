package logrussqlitehook

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type Hook interface {
	logrus.Hook
	Close()
}

func New(path string) (Hook, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS logs (
	level TEXT NOT NULL,
	time TEXT NOT NULL,
	message TEXT NOT NULL,
	data BLOB NOT NULL
);
	`)
	if err != nil {
		return nil, err
	}

	return sqliteHook{
		db: db,
	}, nil
}

type sqliteHook struct {
	db *sql.DB
}

func (h sqliteHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h sqliteHook) Fire(entry *logrus.Entry) error {
	data, err := json.Marshal(entry.Data)
	if err != nil {
		return err
	}

	_, err = h.db.Exec(
		`INSERT INTO logs (level, time, message, data) VALUES (?, ?, ?, ?)`,
		entry.Level.String(),
		entry.Time.Format(time.RFC3339Nano),
		entry.Message,
		data,
	)

	return err
}

func (h sqliteHook) Close() {
	h.db.Close()
}
