package util

import (
	"errors"
	"math"
	"net/http"
	"sync"
	"sync/atomic"
)


func PromiseAllLimiter[T any, R any](data []T, target func(T) R, limit int) ([]R, error) {
	Assert(limit < len(data), "UNKNOWN", http.StatusInternalServerError, "Limit is invalid!")

	if limit <= 0 {
		limit = 1
	}

	var err error
	targetResponse := []R{}

	size := len(data)
	iterations := math.Ceil(float64(size) / float64(limit))

	for i := 0; i < int(iterations); i++ {
		startPos := i * limit
		endPos := (i + 1) * limit
		endPos = Ternary(size < endPos, size, endPos)

		response, e := PromiseAll(data[startPos:endPos], target)

		if e != nil {
			err = e
			break
		} else {
			targetResponse = append(targetResponse, response...)
		}
	}

	if err != nil {
		return nil, err
	} else {
		return targetResponse, err
	}
}

func PromiseAll[T any, R any](dataModels []T, target func(T) R) ([]R, error) {
	var goroutineError atomic.Value
	var responseMap SyncMap
	var wg sync.WaitGroup

	dataLength := len(dataModels)
	wg.Add(dataLength)

	stopCh := make(chan bool)
	for i, data := range dataModels {
		go func(index int, d T) {
			defer errorHandler(&goroutineError, &stopCh, &wg)
			response := target(d)
			responseMap.Store(index, response)
		}(i, data)
	}

	waitForCompletion(&wg, stopCh)

	goroutineErrorValue := goroutineError.Load()
	if goroutineErrorValue == nil {
		targetResponse := orderedResponse[R](&responseMap)
		return targetResponse, nil
	} else {
		return nil, errors.New(goroutineErrorValue.(string))
	}
}

func orderedResponse[R any](responseMap *SyncMap) []R {
	size := responseMap.Size()
	targetResponse := make([]R, size)

	responseMap.Range(func(key any, value any) bool {
		keyVal := key.(int)
		targetResponse[keyVal] = value.(R)

		return true
	})

	return targetResponse
}

func waitForCompletion(wg *sync.WaitGroup, stopCh chan bool) {
	clearCh := make(chan bool)
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(clearCh)
	}(wg)

	go func(wg *sync.WaitGroup) {
		<-stopCh
		close(clearCh)
	}(wg)

	<-clearCh
}

func errorHandler(goroutineError *atomic.Value, stopCh *chan bool, wg *sync.WaitGroup) {
	err := recover()
	if err != nil {
		goroutineError.Store(GetMsgFromError(err))
		close(*stopCh)
	}
	wg.Done()
}

