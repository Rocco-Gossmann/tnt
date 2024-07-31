package database

import (
	"database/sql"
	"log"
	"os"

	sqlite "github.com/rocco-gossmann/go_sqliteutils"
	utils "github.com/rocco-gossmann/go_utils"
)

const DEFAULT_DBFILE string = "tnt.db"

var dbFileName string = DEFAULT_DBFILE

var db *sql.DB

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

	sqlite.InitDBFile(dbFile, 1, func(db *sql.DB, isVersion uint, shouldBeVersion uint) {
		log.Printf("dbfile '%s' does not exist => creating it now\n", dbFile)
		switch isVersion {
		case 0:
			initTasksTable()
			initTimesTable()
		}
	})

	log.Println("DB Init")
}

// Closes the database connection (call via defer, right after you called InitDB)
func DeInitDB() {
	sqlite.DeInitDB()
}

// Check if an error has todo with a unique key already existing
func IsUniqueContraintError(err error) bool {
	return sqlite.IsUniqueContraintError(err)
}
