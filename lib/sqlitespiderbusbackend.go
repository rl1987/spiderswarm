package spiderswarm

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SQLiteSpiderBusBackend struct {
	SpiderBusBackend

	dbConn         *sql.DB
	sqliteFilePath string
}

func NewSQLiteSpiderBusBackend(sqliteFilePath string) *SQLiteSpiderBusBackend {
	if sqliteFilePath == "" {
		sqliteDirPath, err := ioutil.TempDir(os.TempDir(), "spiderbus_")
		if err != nil {
			log.Error(fmt.Sprintf("Failed to create temp dir for SQLiteSpiderBusBackend: %v", err))
			return nil
		}

		sqliteFilePath = sqliteDirPath + "/" + "spiderbus.db"
	}

	dbConn, err := sql.Open("sqlite3", sqliteFilePath)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to open DB connection: %v", err))
		return nil
	}

	tx, _ := dbConn.Begin()

	tx.Exec("CREATE TABLE IF NOT EXISTS scheduledTasks (id INTEGER PRIMARY KEY AUTOINCREMENT, raw BLOB)")
	tx.Exec("CREATE TABLE IF NOT EXISTS taskPromises (id INTEGER PRIMARY KEY AUTOINCREMENT, raw BLOB)")
	tx.Exec("CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY AUTOINCREMENT, raw BLOB)")

	tx.Commit()

	fmt.Printf("Created new SQLite DB at: %s\n", sqliteFilePath)

	return &SQLiteSpiderBusBackend{
		dbConn:         dbConn,
		sqliteFilePath: sqliteFilePath,
	}
}

func (ssbb *SQLiteSpiderBusBackend) encodeEntry(entry interface{}) []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)

	encoder.Encode(entry)

	bytes, _ := ioutil.ReadAll(buffer)

	return bytes
}

func (ssbb *SQLiteSpiderBusBackend) decodeEntry(raw []byte, entry interface{}) interface{} {
	buffer := bytes.NewBuffer(raw)
	decoder := json.NewDecoder(buffer)

	err := decoder.Decode(entry)
	if err != nil {
		spew.Dump(err)
	}

	return &entry
}

func (ssbb *SQLiteSpiderBusBackend) SendScheduledTask(scheduledTask *ScheduledTask) error {
	raw := ssbb.encodeEntry(scheduledTask)

	tx, _ := ssbb.dbConn.Begin()

	tx.Exec("INSERT INTO scheduledTasks (raw) VALUES (?)", raw)

	tx.Commit()

	return nil
}

func (ssbb *SQLiteSpiderBusBackend) ReceiveScheduledTask() *ScheduledTask {
	tx, _ := ssbb.dbConn.Begin()

	var row_id int
	var raw []byte

	row := tx.QueryRow("SELECT * FROM scheduledTasks ORDER BY id ASC LIMIT 1")

	err := row.Scan(&row_id, &raw)
	if err != nil {
		spew.Dump(err)
		return nil
	}

	scheduledTask := &ScheduledTask{}

	ssbb.decodeEntry(raw, scheduledTask)

	tx.Exec(fmt.Sprintf("DELETE FROM scheduledTasks WHERE id=%d", row_id))
	tx.Commit()

	return scheduledTask
}

func SendTaskPromise(taskPromise *TaskPromise) error {
	return errors.New("Not implemented")
}

func ReceiveTaskPromise() *TaskPromise {
	return nil
}

func SendItem(item *Item) error {
	return errors.New("Not implemented")
}

func ReceiveItem() *Item {
	return nil
}
