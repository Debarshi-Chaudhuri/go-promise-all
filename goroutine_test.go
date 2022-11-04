package util

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)


func TestFailure(t *testing.T) {
	ids := []time.Duration{12, 2, 23, 34, 5, 6, 73, 8, 19}

	timeNow1 := time.Now()
	vals, err := PromiseAll(ids, func(id time.Duration) time.Duration {
		time.Sleep(time.Second * id)
		if id == time.Duration(5) {
			Assert(false, "UNKNOWN", http.StatusInternalServerError, "Time Exceeded!!")
		}

		return id
	})
	timeNow2 := time.Now()
	timeDiff := timeNow2.Sub(timeNow1) / time.Second

	fmt.Println(vals, err, timeDiff)
}

func TestSuccess(t *testing.T) {
	ids := []time.Duration{11, 2, 25, 4, 15, 6, 17, 8, 19}

	timeNow1 := time.Now()
	vals, err := PromiseAll(ids, func(id time.Duration) time.Duration {
		time.Sleep(time.Second * id)

		return id
	})
	timeNow2 := time.Now()
	timeDiff := timeNow2.Sub(timeNow1) / time.Second

	fmt.Println(vals, err, timeDiff)
}

func TestFailureLimiter(t *testing.T) {
	ids := []time.Duration{5, 4, 3, 2, 1, 5, 3, 4, 5, 1}

	timeNow1 := time.Now()
	vals, err := PromiseAllLimiter(ids, func(id time.Duration) time.Duration {
		time.Sleep(time.Second * id)

		if id == time.Duration(1) {
			Assert(false, "UNKNOWN", http.StatusInternalServerError, "Time Exceeded!!")
		}

		return id
	}, 4)
	timeNow2 := time.Now()
	timeDiff := timeNow2.Sub(timeNow1) / time.Second

	fmt.Println(vals, err, timeDiff)
}

func TestSuccessLimiter(t *testing.T) {
	ids := []time.Duration{5, 4, 3, 2, 1, 5, 3}

	timeNow1 := time.Now()
	vals, err := PromiseAllLimiter(ids, func(id time.Duration) time.Duration {
		time.Sleep(time.Second * id)

		return id
	}, 4)
	timeNow2 := time.Now()
	timeDiff := timeNow2.Sub(timeNow1) / time.Second

	fmt.Println(vals, err, timeDiff)
}
