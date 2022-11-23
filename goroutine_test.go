package util

import (
	"context"
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


func TestSuccessExtended(t *testing.T) {
	ctx := context.Background()

	var time1 time.Duration = 3
	var time2 time.Duration = 4
	var time3 time.Duration = 5

	ctx = context.WithValue(ctx, 1, time1)
	ctx = context.WithValue(ctx, 2, time2)
	ctx = context.WithValue(ctx, 3, time3)
	timeNow1 := time.Now()

	val := uint(1)
	p1 := Promisify(f1)
	p2 := Promisify(f2, &val)
	p3 := Promisify(f3)

	vals, err := PromiseAllExtended(&ctx, p1, p2, p3)
	timeNow2 := time.Now()
	timeDiff := timeNow2.Sub(timeNow1) / time.Second

	fmt.Println(vals, err, timeDiff)
}

func TestFailedExtended(t *testing.T) {
	ctx := context.Background()

	var time1 time.Duration = 3
	var time2 time.Duration = 4
	var time3 time.Duration = 5

	ctx = context.WithValue(ctx, 1, time1)
	ctx = context.WithValue(ctx, 2, time2)
	ctx = context.WithValue(ctx, 3, time3)
	timeNow1 := time.Now()

	val := uint(1)
	p1 := Promisify(f1)
	p2 := Promisify(f2, &val)
	p3 := Promisify(f4)

	vals, err := PromiseAllExtended(&ctx, p1, p2, p3)
	timeNow2 := time.Now()
	timeDiff := timeNow2.Sub(timeNow1) / time.Second

	fmt.Println(vals, err, timeDiff)
}

func f1(ctx *context.Context) interface{} {
	val := (*ctx).Value(1).(time.Duration)
	time.Sleep(val * time.Second)
	return val
}

func f2(ctx *context.Context, duration *uint) *uint {
	val := (*ctx).Value(2).(time.Duration)
	time.Sleep(val * time.Second)
	return duration
}

func f3(ctx *context.Context) interface{} {
	val := (*ctx).Value(3).(time.Duration)
	time.Sleep(val * time.Second)
	return 10
}

func f4(ctx *context.Context) interface{} {
	val := (*ctx).Value(1).(time.Duration)
	time.Sleep(val * time.Second)
	Assert(false, "UNKNOWN", http.StatusInternalServerError, "Time Exceeded!!")
	
	return val
}
