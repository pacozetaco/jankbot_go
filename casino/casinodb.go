package casino

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	//init all the tables if not exist on startup
	err := createDBConnection()
	if err != nil {
		log.Fatal("error connecting to db", err)
	}

	defer closeDBConnection()

	createJankCoinsTable()
	createBJTable()
	createHiLoTable()
	createDeathrollTable()

}

// func for all db connections
func createDBConnection() error {
	user := os.Getenv("SQL_USER")
	pass := os.Getenv("SQL_PASS")
	host := os.Getenv("SQL_HOST")
	port := os.Getenv("SQL_PORT")
	dbName := os.Getenv("SQL_DB")
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbName)
	var err error
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		return err
	}
	return nil
}

// dunc to close the conection
func closeDBConnection() {
	err := db.Close()
	if err != nil {
		log.Println("error closing db connection", err)
	}
}

func createJankCoinsTable() {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS jankcoins (
                name VARCHAR(255) PRIMARY KEY UNIQUE,
                coins BIGINT, 
                lastclaim DATE
    );
    `)
	if err != nil {
		log.Println("error creating jankcoins table", err)
	} else {
		log.Println("jankcoins table created")
	}
}

func createHiLoTable() {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS hilo_log (
            id BIGINT AUTO_INCREMENT PRIMARY KEY,
            date DATE,
            time TIME,
            player TEXT,
            bet BIGINT,
            choice TEXT,
            roll INT,
            result TEXT
    );
    `)
	if err != nil {
		log.Println("error creating hilo table", err)
	} else {
		log.Println("hilo table created")
	}
}

func createBJTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS bj_log (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			date DATE,
			time TIME,
			player TEXT,
			bet BIGINT,
			result TEXT
	);
	`)
	if err != nil {
		log.Println("error creating bj table", err)
	} else {
		log.Println("bj table created")
	}
}

func createDeathrollTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS deathroll_log (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			date DATE,
			time TIME,
			player TEXT,
			bet BIGINT,
			whofirst TEXT,
			result TEXT,
			gamecontent TEXT
	);
	`)
	if err != nil {
		log.Println("error creating deathroll table", err)
	} else {
		log.Println("deathroll table created")
	}
}

// get coin balance
func getBalance(name string) (int, error) {
	err := createDBConnection()
	if err != nil {
		return -2, err
	}
	defer closeDBConnection()

	var balance int
	err = db.QueryRow("SELECT coins FROM jankcoins WHERE name = ?", name).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.New("balance not found")
		}
		return -2, err
	}
	return balance, nil
}

// add or subtract coins
func addBalance(name string, coins int) error {
	err := createDBConnection()
	if err != nil {
		return err
	}

	defer closeDBConnection()

	_, err = db.Exec("UPDATE jankcoins SET coins = coins + ? WHERE name = ?", coins, name)
	return err
}

// daily coins, set at 100 rn
func dailyCoins(name string) string {

	var lastClaim string
	var coins int

	now := time.Now()

	err := createDBConnection()
	if err != nil {
		return "1 error connecting to db"
	}

	defer closeDBConnection()

	row := db.QueryRow("SELECT coins, lastclaim FROM jankcoins WHERE name = ?", name)

	switch err := row.Scan(&coins, &lastClaim); err {
	case sql.ErrNoRows:
		_, err = db.Exec("INSERT INTO jankcoins (name, coins, lastclaim) VALUES (?, ?, ?)", name, 100, now)
		if err != nil {
			return "error connecting to db"
		}
		return "100 coins added!"
	case nil:
		claimTime, err := time.Parse("2006-01-02", lastClaim)
		if err != nil {
			return "error parsing last claim time"
		}
		if claimTime.Format("2006-01-02") == now.Format("2006-01-02") {
			bal := strconv.Itoa(coins)
			return "You already claimed today! Balance: " + bal
		} else {
			_, err = db.Exec("UPDATE jankcoins SET coins = ?, lastclaim = ? WHERE name = ?", coins+100, now, name)
			if err != nil {
				return "error connecting to db"
			}
			return "100 coins added! Your new balance is: " + fmt.Sprint(coins+100) + " coins"

		}
	default:
		panic(err)
	}
}

// log hilo game
func (h *hiLoG) logHiLo() {
	createDBConnection()
	defer closeDBConnection()
	stmt, err := db.Prepare("INSERT INTO hilo_log (date, time, player, bet, choice, roll, result) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(time.Now().Format("2006-01-02"), time.Now().Format("15:04:05"), h.player, h.bet, h.choice, h.roll, h.result)
	if err != nil {
		return
	}
}

func (d *deathRollG) logDeathRoll() {
	createDBConnection()
	defer closeDBConnection()
	stmt, err := db.Prepare("INSERT INTO deathroll_log (date, time, player, bet, whofirst, result, gamecontent) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(time.Now().Format("2006-01-02"), time.Now().Format("15:04:05"), d.player, d.bet, d.first, d.result, d.msg.Content)
	if err != nil {
		return
	}
}

func (b *blackJackG) logBJ() {
	createDBConnection()
	defer closeDBConnection()
	stmt, err := db.Prepare("INSERT INTO bj_log (date, time, player, bet, result) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(time.Now().Format("2006-01-02"), time.Now().Format("15:04:05"), b.player, b.bet, b.result)
	if err != nil {
		return
	}
}

//TODO

// GAME LOGS - STATS
