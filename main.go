package main

import (
	"context"
	"encoding/json"
	"fmt"
	containerlib "github.com/resource-aware-jds/container-lib"
	"github.com/resource-aware-jds/container-lib/model"
	"github.com/resource-aware-jds/container-lib/pkg/containerlibcontext"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TaskAttribute struct {
	SleepTime            *int `json:"sleepTime,omitempty"`
	MemoryAllocationSize *int `json:"memoryAllocationSize,omitempty"`
}

func main() {
	containerlib.Run(func(ctx containerlibcontext.Context, task model.Task) error {
		var unmarshalledData TaskAttribute
		err := json.Unmarshal(task.Attributes, &unmarshalledData)
		if err != nil {
			logrus.Error(err)
			return err
		}

		sleepTime := 15 * time.Second
		if unmarshalledData.SleepTime != nil {
			sleepTime = time.Duration(*unmarshalledData.SleepTime) * time.Second
		}

		if unmarshalledData.MemoryAllocationSize == nil {
			size := 0
			unmarshalledData.MemoryAllocationSize = &size
		}

		innerCtx := context.Background()
		newCtx, cancelFunc := context.WithCancel(innerCtx)
		a := make([]byte, (*unmarshalledData.MemoryAllocationSize)*1000000)
		fmt.Println("Reserved: ", " : ", *unmarshalledData.MemoryAllocationSize, "mb")
		fmt.Println("Press Ctrl+C to abort this operation")
		go func(innerCtx context.Context) {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Done.")
					return
				default:
					for i := range a {
						a[i] = 0x99
					}
				}
				fmt.Println("Sleep 10 Second before reserving the next")
				time.Sleep(10 * time.Second)
			}

		}(newCtx)

		// Gracefully Shutdown
		// Make channel listen for signals from OS
		gracefulStop := make(chan os.Signal, 1)
		signal.Notify(gracefulStop, syscall.SIGTERM)
		signal.Notify(gracefulStop, syscall.SIGINT)

		<-gracefulStop

		cancelFunc()

		time.Sleep(sleepTime)
		fmt.Println(unmarshalledData)

		ctx.Success()
		ctx.RecordResult(task.Attributes)
		return nil
	})
}
