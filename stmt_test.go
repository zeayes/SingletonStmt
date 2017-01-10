package SingletonStmt

import (
	"database/sql"
	"log"
	"strconv"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	BaseSQL = "SELECT id FROM test WHERE id "
)

var (
	db        *sql.DB
	singleton *SingletonStmt
	IDS       = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
)

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/information_schema?timeout=3s")
	if err != nil {
		log.Fatalf("open db error: %v", err)
	}
	if err = setup(db); err != nil {
		log.Fatalf("init table error: %v", err)
	}
	singleton = NewSingletonStmt(db, BaseSQL+" = ?")
}

func setup(db *sql.DB) (err error) {
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS `test`")
	if err != nil {
		return
	}
	_, err = db.Exec("USE `test`")
	if err != nil {
		return
	}
	_, err = db.Exec("DROP TABLE IF EXISTS `test`")
	if err != nil {
		return
	}
	_, err = db.Exec("CREATE TABLE test (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(10))")
	if err != nil {
		return
	}
	_, err = db.Exec(`INSERT INTO test (name) VALUES ('test'), ('test'), ('test'),
			('test'), ('test'), ('test'), ('test'), ('test'), ('test'), ('test')`)
	if err != nil {
		return
	}
	return
}

func TestSingletonStmt(t *testing.T) {
	if err := singleton.GetStmt(); err != nil {
		t.Fatalf("TestSingletonStmt GetStmt error: %v", err)
	}
	expect := 1
	row := singleton.Stmt.QueryRow(expect)
	var v sql.NullInt64
	if err := row.Scan(&v); err != nil {
		t.Fatalf("TestSingletonStmt Scan error: %v", err)
	}
	if v.Int64 != int64(expect) {
		t.Fatalf("TestSingletonStmt result expect %v but got %v", expect, v.Int64)
	}
}

func BenchmarkDefaultQuery(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		row := db.QueryRow(BaseSQL + " = " + strconv.Itoa(1))
		var v sql.NullInt64
		if err := row.Scan(&v); err != nil {
			b.Errorf("BenchmarkDefaultQuery Scan error: %v", err)
		}
	}
}

func BenchmarkDefaultStmt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stmt, err := db.Prepare(BaseSQL + " = ?")
		if err != nil {
			b.Errorf("BenchmarkDefaultStmt Prepare error: %v", err)
		}
		row := stmt.QueryRow(1)
		var v sql.NullInt64
		if err = row.Scan(&v); err != nil {
			b.Errorf("BenchmarkDefaultStmt Scan error: %v", err)
		}
		stmt.Close()
	}
}

func BenchmarkSingletonStmt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	if err := singleton.GetStmt(); err != nil {
		b.Errorf("BenchmarkSingletonStmt GetStmt error: %v", err)
	}
	for i := 0; i < b.N; i++ {
		row := singleton.Stmt.QueryRow(1)
		var v sql.NullInt64
		if err := row.Scan(&v); err != nil {
			b.Errorf("BenchmarkSingletonStmt Scan error: %v", err)
		}
	}
}
