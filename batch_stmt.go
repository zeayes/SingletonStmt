package SingletonStmt

import (
	"database/sql"
	"strings"
	"sync"
)

// BatchSingletonStmt used for `SELECT columns FROM table WHERE id IN (1, 2, 3)`
type BatchSingletonStmt struct {
	sync.Mutex
	db      *sql.DB
	baseSQL string
	stmts   []*SingletonStmt
}

// NewBatchSingletonStmt create new BatchSingletonStmt
func NewBatchSingletonStmt(db *sql.DB, baseSQL string, count int) *BatchSingletonStmt {
	return &BatchSingletonStmt{
		db:      db,
		baseSQL: baseSQL,
		stmts:   make([]*SingletonStmt, count),
	}
}

// GetStmt get a stmt with index, index means `len(ids) - 1`
func (bs *BatchSingletonStmt) GetStmt(index int) (*sql.Stmt, error) {
	stmt := bs.stmts[index]
	if stmt != nil {
		return stmt.Stmt, nil
	}
	bs.Lock()
	stmt = bs.stmts[index]
	if stmt != nil {
		bs.Unlock()
		return stmt.Stmt, nil
	}
	clause := strings.TrimRight(strings.Repeat("?,", index+1), ",")
	ss := NewSingletonStmt(bs.db, bs.baseSQL+" IN ("+clause+")")
	if err := ss.GetStmt(); err != nil {
		bs.Unlock()
		return nil, err
	}
	bs.stmts[index] = ss
	bs.Unlock()
	return ss.Stmt, nil
}
