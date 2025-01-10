package main

import (
	// "net/http"
	// _ "net/http/pprof"

	"github.com/rlawnsxo131/ws-placeholder/api"
)

func main() {
	// go func() {
	// 	http.ListenAndServe("0.0.0.0:6060", nil)
	// }()

	api.Run("8080")
}
