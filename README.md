# performance

- [pprof graph 참고문서](https://github.com/google/pprof/blob/main/doc/README.md#interpreting-the-callgraph)

```sh
# 시각화를 위한 의존성 설치
$ brew install graphviz
```

```go
// 프로파일링을 위한 port 추가
import (
    ...

	_ "net/http/pprof"

    ...
)

func main() {
    ...

    go func() {
     http.ListenAndServe("0.0.0.0:6060", nil)
    }()

    ...
}
```

```sh
# 테스트를 위한 빌드및 애플리케이션 실행
$ go build ./cmd/api/main.go
$ ./main

# 트래픽 넣기
# request count: 10000
# concurrency count: 100
# max test time: 6000sec
# keep-alive: true
$ ab -n 10000 -c 100 -t 6000 -k http://localhost:8080/test
```

```sh
# pprof 보기
$ go tool pprof -http 0.0.0.0:3000 http://0.0.0.0:6060/debug/pprof/heap

# allocs: A sampling of all past memory allocations
# block: Stack traces that led to blocking on synchronization primitives
# cmdline: The command line invocation of the current program
# goroutine: Stack traces of all current goroutines. Use debug=2 as a query # parameter to export in the same format as an unrecovered panic.
# heap: A sampling of memory allocations of live objects. You can specify the gc GET # parameter to run GC before taking the heap sample.
# mutex: Stack traces of holders of contended mutexes
# profile: CPU profile. You can specify the duration in the seconds GET parameter. # After you get the profile file, use the go tool pprof command to investigate the # profile.
# threadcreate: Stack traces that led to the creation of new OS threads
# trace: A trace of execution of the current program. You can specify the duration # in the seconds GET parameter. After you get the trace file, use the go tool trace # command to investigate the trace.

# trace
$ curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
$ go tool trace trace.out

# profile
$ curl -o profile.out http://0.0.0.0:6060/debug/pprof/profile\?seconds\=180
$ go tool pprof -http 0.0.0.0:3000 profile.out
```
