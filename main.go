package main

import (
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nithishravindra/conn-pool-sql/internal/mysql"
)

var SQLQUERY = "SELECT SLEEP(0.1);"

func withoutPool(dbConfig mysql.MySQLConfig, noQuery int) {
	startTime := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < noQuery; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db, err := mysql.GetNewConnection(dbConfig)
			if err != nil {
				panic(err)
			}
			_, err = db.Exec(SQLQUERY)
			if err != nil {
				panic(err)
			}
			db.Close()
		}()
	}
	wg.Wait()
	log.Println("withoutPool:", time.Since(startTime))
}

func withPool(dbConfig mysql.MySQLConfig, noQuery int) {
	startTime := time.Now()

	var wg sync.WaitGroup
	pool, err := mysql.GetConnPool(dbConfig, 150)
	if err != nil {
		panic(err)
	}

	for i := 0; i < noQuery; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := pool.Get()
			if err != nil {
				panic(err)
			}
			_, err = conn.Db.Exec(SQLQUERY)
			if err != nil {
				panic(err)
			}
			pool.Put(conn)
		}()
	}

	log.Println("withPool:", time.Since(startTime))
}

func main() {
	config := mysql.MySQLConfig{
		UserName: "pool",
		Password: "pool",
		Port:     3306,
		DbName:   "connpool",
	}
	withoutPool(config, 100)
	withPool(config, 80000)
}
