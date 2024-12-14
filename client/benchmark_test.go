package main

import (
    "client/original"
    "client/server"
    "testing"
)

var pairs = []string{"BTCUSDT", "ETHUSDT", "BNBUSDT", "ADAUSDT", "XRPUSDT"}
const timeoutBenchmark = 10

// run with go test -bench=.
func BenchmarkOriginalInit(b *testing.B) {
    for i := 0; i < b.N; i++ {
        original.Init(pairs, false, timeoutBenchmark)
    }
}

func BenchmarkServerInit(b *testing.B) {
    for i := 0; i < b.N; i++ {
        server.Init(pairs, false, timeoutBenchmark)
    }
}