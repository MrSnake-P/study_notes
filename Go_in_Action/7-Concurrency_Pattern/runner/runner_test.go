package runner

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestRunnerTimeout(t *testing.T) {
	const timeout = 3 * time.Second

	r := New(timeout)
	r.Add(creatTask(), creatTask(), creatTask())
	err := r.Start()

	require.Equal(t, ErrTimeout, err)
}

func creatTask() func(int) {
	return func(id int) {
		log.Printf("Processor - Task #%d.", id)
		time.Sleep(time.Duration(id) * time.Second)
	}
}
