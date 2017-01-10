package SingletonStmt

import (
	"database/sql"
	"sync"
)

// SingletonStmt sql.Stmt singleton
type SingletonStmt struct {
	sync.Once
	db    *sql.DB
	query string
	Stmt  *sql.Stmt
}

// NewSingletonStmt create a new SingletonStmt instance
func NewSingletonStmt(db *sql.DB, query string) *SingletonStmt {
	return &SingletonStmt{db: db, query: query}
}

// GetStmt get a prepare stmt by SQL
func (ss *SingletonStmt) GetStmt() (err error) {
	if ss.Stmt != nil {
		return
	}
	ss.Do(func() {
		ss.Stmt, err = ss.db.Prepare(ss.query)
		if err != nil {
			return
		}
	})
	return
}
