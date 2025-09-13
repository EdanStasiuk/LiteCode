package sandbox

import (
	"math/rand"
	"time"
)

func RunCode(code string, language string) (status string, result string, runtime float64, memory int64) {
	// TODO: Run code in Docker sandbox, capture output, runtime, memory
	time.Sleep(time.Millisecond * 500) // simulate execution

	// Random stub results for testing
	status = "success"
	result = "Accepted"
	runtime = rand.Float64() * 2
	memory = rand.Int63n(100_000_000)
	return
}
