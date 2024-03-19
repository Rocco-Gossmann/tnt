package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rocco-gossmann/tnt/pkg/utils"
)

const DEFAULT_DBFILE string = "tnt.db"

var dbFileName string = DEFAULT_DBFILE

var db *sql.DB

// Runs a prepared statement on the database. Requires the DB to be initilaized first
func ExecStatement(statement string, args ...any) (sql.Result, error) {
	stmt, err := db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

// Runs a prepared statement on the database. Requires the DB to be initilaized first
func QueryStatement(statement string, args ...any) (*sql.Rows, error) {
	stmt, err := db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Query(args...)
}

// Runs a prepared statement on the database. Requires the DB to be initilaized first
func RowQueryStatement(statement string, args ...any) (*sql.Row, error) {
	stmt, err := db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}

// Sets a new Default db file path
func SetDBFileName(fileName string) {
	dbFileName = fileName
}

// Opens or creates the database
func InitDB(dbFile string) {
	var err error

	if db != nil {
		log.Println("DB Open")
		return
	}

	if len(dbFile) == 0 {
		dbFile = dbFileName
	}

	if len(dbFile) == 0 {

		configDir, err := os.UserHomeDir()
		utils.Err(err)

		configPath := configDir + "/.local/share/tnt"
		_, err = os.Stat(configPath)
		if os.IsNotExist(err) {
			// create Directory if not exist
			err := os.MkdirAll(configPath, 0700)
			utils.Err(err)
		}

		dbFile = configPath + "/" + DEFAULT_DBFILE
	}

	utils.Err(err)

	_, fiErr := os.Stat(dbFile)
	db, err = sql.Open("sqlite3", dbFile)
	if os.IsNotExist(fiErr) {
		log.Printf("dbfile '%s' does not exist => creating it now\n", dbFile)
		initTasksTable()
		initTimesTable()
	}

	log.Println("DB Init")
}

// Closes the database connection (call via defer, right after you called InitDB)
func DeInitDB() {
	if db != nil {
		db.Close()
		db = nil
		log.Println("DB DeInit")
	}

}

// Check if an error has todo with a unique key already existing
func IsUniqueContraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasPrefix(err.Error(), "UNIQUE constraint failed")
}
