package main

import (
	"encoding/json"
	"fmt"
	containerlib "github.com/resource-aware-jds/container-lib"
	"github.com/resource-aware-jds/container-lib/model"
	"github.com/resource-aware-jds/container-lib/pkg/containerlibcontext"
	"github.com/sirupsen/logrus"
	"time"
)

type TaskAttribute struct {
	SleepTime            *time.Duration `json:"sleepTime,omitempty"`
	MemoryAllocationSize *int           `json:"memoryAllocationSize,omitempty"`
}

func main() {
	containerlib.Run(func(ctx containerlibcontext.Context, task model.Task) error {
		var unmarshalledData TaskAttribute
		err := json.Unmarshal(task.Attributes, &unmarshalledData)
		if err != nil {
			logrus.Error(err)
			return err
		}

		if unmarshalledData.SleepTime == nil {
			sleepTime := 15 * time.Second
			unmarshalledData.SleepTime = &sleepTime
		}

		if unmarshalledData.MemoryAllocationSize == nil {
			size := 0
			unmarshalledData.MemoryAllocationSize = &size
		}

		fakeAllocation := make([]byte, *unmarshalledData.MemoryAllocationSize)
		for i := 0; i < *unmarshalledData.MemoryAllocationSize; i++ {
			fakeAllocation[i] = 0x99
		}

		time.Sleep(*unmarshalledData.SleepTime)
		fmt.Println(unmarshalledData)

		ctx.Success()
		ctx.RecordResult(task.Attributes)
		return nil
	})
}
