package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func WaitInitialization(port string) {

	db, err := sql.Open("mysql", "root:@(127.0.0.1:"+port+")/")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var retryCnt int

	time.Sleep(4 * time.Second)
	for {
		var initialized int
		if err := db.QueryRow("SELECT GET_LOCK('initialize', -1) AS initialized").Scan(&initialized); err != nil {
			if retryCnt > 5 {
				panic(err.Error())
			}
		}

		if initialized == 1 {
			break
		}
		if retryCnt > 10 {
			panic("can not get lock.")
		}
		retryCnt++
		time.Sleep(1 * time.Second)
	}
}
