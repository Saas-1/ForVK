package ForVK

import (
	"fmt"
	"time"
)

func main() {
	pool := CreateWorkerPool()

	pool.AddWorker()
	pool.AddWorker()

	time.Sleep(time.Second)

	for i := 1; i <= 10; i++ {
		pool.AddJob(fmt.Sprintf("Job %d", i))
	}

	time.Sleep(time.Second)

	pool.RemoveWorker()

	time.Sleep(time.Second)

	for i := 6; i <= 10; i++ {
		pool.AddJob(fmt.Sprintf("Job %d", i))
	}

	time.Sleep(time.Second)
}
