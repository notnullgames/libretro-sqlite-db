package main

import (
	"fmt"
	"io/ioutil"
	"database/sql"
	"strconv"
	"strings"
	"github.com/libretro/ludo/rdb"
	_ "github.com/mattn/go-sqlite3"
)

type RDBEntry struct {
	filename string
	system string
}

func main() {
	var entries = []RDBEntry {
		RDBEntry {
			filename: "libretro-database/rdb/Nintendo - Game Boy.rdb",
			system: "gb",
		},
		RDBEntry {
			filename: "libretro-database/rdb/Nintendo - Game Boy Color.rdb",
			system: "gbc",
		},
		RDBEntry {
			filename: "libretro-database/rdb/Nintendo - Game Boy Advance.rdb",
			system: "gba",
		},
		RDBEntry {
			filename: "libretro-database/rdb/Nintendo - Nintendo 64.rdb",
			system: "n64",
		},
		RDBEntry {
			filename: "libretro-database/rdb/Sega - Mega Drive - Genesis.rdb",
			system: "md",
		},
		RDBEntry {
			filename: "libretro-database/rdb/Nintendo - Nintendo Entertainment System.rdb",
			system: "nes",
		},
		RDBEntry {
			filename: "libretro-database/rdb/Nintendo - Super Nintendo Entertainment System.rdb",
			system: "snes",
		},
	}

	database, _ := sql.Open("sqlite3", "./libretro-db.sqlite")

	exec(database, "CREATE TABLE IF NOT EXISTS games (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT, system TEXT, romName TEXT, developer TEXT, crc32 TEXT)")

	for _, entry := range entries {
		loadDatabase(database, entry.filename, entry.system)
	}

	exec(database, "CREATE INDEX crc32Index ON games (crc32)")
	exec(database, "CREATE INDEX romNameIndex ON games (romName)")
}

func loadDatabase(database *sql.DB, filename string, system string) {
	bytes, _ := ioutil.ReadFile(filename)
	var games = rdb.Parse(bytes)

	for i, g := range games {
		fmt.Println(i, g.Name)
		crc32 := strings.ToUpper(strconv.FormatInt(int64(g.CRC32), 16))
		statement, _ := database.Prepare("INSERT INTO games (name, romName, system, developer, crc32) VALUES (?,?,?,?,?)")
		statement.Exec(g.Name, g.ROMName, system, g.Developer, crc32)
	}
}

func exec(database *sql.DB, sqlString string) {
	statement, _ := database.Prepare(sqlString)
	statement.Exec()
}