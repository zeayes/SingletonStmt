SingletonStmt
===========
[![Build Status](https://travis-ci.org/zeayes/SingletonStmt.svg?branch=master)](https://travis-ci.org/zeayes/SingletonStmt)

A goroutine safe sql.Stmt for global usage in your application.

Install
===========
```bash
go get -u github.com/zeayes/SingletonStmt
```

Benchmark
===========
benchmark on MBP(Mid 2015 2.2 GHz 16GB), and memcached served by default options.
```
go test -run=^$ -bench .
BenchmarkDefaultBatchQuery-8    	    2000	    996211 ns/op	    6593 B/op	     400 allocs/op
BenchmarkDefaultBatchStmt-8     	    1000	   1629131 ns/op	   14369 B/op	     431 allocs/op
BenchmarkBatchSingletonStmt-8   	    2000	    898429 ns/op	    7002 B/op	     320 allocs/op
BenchmarkDefaultQuery-8         	   20000	     80885 ns/op	     372 B/op	      14 allocs/op
BenchmarkDefaultStmt-8          	   10000	    142168 ns/op	     822 B/op	      21 allocs/op
BenchmarkSingletonStmt-8        	   20000	     76495 ns/op	     312 B/op	      13 allocs/op
```

License
===========
*SingletonStmt* is released under the MIT License.
