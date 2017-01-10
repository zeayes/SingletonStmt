package SingletonStmt

import (
	"database/sql"
	"strconv"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestBatchSingletonStmt(t *testing.T) {
	batchStmt := NewBatchSingletonStmt(db, BaseSQL, len(IDS))
	for j := 0; j < len(IDS); j++ {
		stmt, err := batchStmt.GetStmt(j)
		if err != nil {
			t.Fatalf("TestBatchSingletonStmt GetStmt error: %v", err)
		}
		args := make([]interface{}, 0, j+1)
		for _, id := range IDS[:j+1] {
			args = append(args, id)
		}
		rows, err := stmt.Query(args...)
		if err != nil {
			t.Fatalf("TestBatchSingletonStmt error: %v", err)
		}
		defer rows.Close()
		for rows.Next() {
			var v sql.NullInt64
			if err = rows.Scan(&v); err != nil {
				t.Fatalf("TestBatchSingletonStmt Scan error: %v", err)
			}
		}
	}
}

func BenchmarkDefaultBatchQuery(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			clause := strconv.Itoa(IDS[0])
			for _, id := range IDS[1 : j+1] {
				clause += ","
				clause += strconv.Itoa(id)
			}
			rows, err := db.Query(BaseSQL + " IN (" + clause + ")")
			if err != nil {
				b.Errorf("BenchmarkDefaultBatchQuery error: %v", err)
			}
			for rows.Next() {
				var v sql.NullInt64
				if err = rows.Scan(&v); err != nil {
					b.Errorf("BenchmarkDefaultBatchQuery Scan error: %v", err)
				}
			}
		}
	}
}

func BenchmarkDefaultBatchStmt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			ids := IDS[:j+1]
			clause := strings.TrimRight(strings.Repeat("?,", j+1), ",")
			stmt, err := db.Prepare(BaseSQL + " IN (" + clause + ")")
			if err != nil {
				b.Errorf("BenchmarkDefaultBatchStmt Prepare error: %v", err)
			}
			args := make([]interface{}, 0, j+1)
			for _, id := range ids {
				args = append(args, id)
			}
			rows, err := stmt.Query(args...)
			if err != nil {
				b.Errorf("BenchmarkDefaultBatchStmt error: %v", err)
			}
			defer rows.Close()
			for rows.Next() {
				var v sql.NullInt64
				if err = rows.Scan(&v); err != nil {
					b.Errorf("BenchmarkDefaultBatchStmt Scan error: %v", err)
				}
			}
		}
	}
}

func BenchmarkBatchSingletonStmt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	batchStmt := NewBatchSingletonStmt(db, BaseSQL, len(IDS))
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			stmt, err := batchStmt.GetStmt(j)
			if err != nil {
				b.Errorf("BenchmarkSingletonBatchStmt GetStmt error: %v", err)
			}
			args := make([]interface{}, 0, j+1)
			for _, id := range IDS[:j+1] {
				args = append(args, id)
			}
			rows, err := stmt.Query(args...)
			if err != nil {
				b.Errorf("BenchmarkSingletonBatchStmt error: %v", err)
			}
			defer rows.Close()
			for rows.Next() {
				var v sql.NullInt64
				if err = rows.Scan(&v); err != nil {
					b.Errorf("BenchmarkSingletonBatchStmt Scan error: %v", err)
				}
			}
		}
	}
}
