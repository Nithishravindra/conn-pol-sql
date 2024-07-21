package mysql

import (
	"database/sql"
	"fmt"
	"sync"
)

type MySQLConfig struct {
	UserName string
	Password string
	Port     int
	DbName   string
}

func GetNewConnection(config MySQLConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s", config.UserName, config.Password, config.Port, config.DbName))
	if err != nil {
		fmt.Println("error opening connection:", err)
		return nil, err
	}
	return db, nil
}

type conn struct {
	Db *sql.DB
}

type ConnPool struct {
	mu      sync.Mutex
	channel chan interface{}
	conz    []*conn
	maxCon  int
}

func GetConnPool(dbConfig MySQLConfig, noOfConn int) (*ConnPool, error) {
	pool := &ConnPool{
		conz:    make([]*conn, 0, noOfConn),
		maxCon:  noOfConn,
		channel: make(chan interface{}, noOfConn),
	}

	for i := 0; i < noOfConn; i++ {
		db, err := GetNewConnection(dbConfig)
		if err != nil {
			return nil, err
		}
		pool.conz = append(pool.conz, &conn{Db: db})
		pool.channel <- nil
	}
	return pool, nil
}

func (pool *ConnPool) Get() (*conn, error) {
	<-pool.channel // blocking..

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if len(pool.conz) == 0 {
		return nil, fmt.Errorf("conn pool exhausted")
	}

	con := pool.conz[0]
	pool.conz = pool.conz[1:]
	return con, nil
}

func (pool *ConnPool) Put(c *conn) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.conz = append(pool.conz, c)
	pool.channel <- nil
}
