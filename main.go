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

func main() {
	containerlib.Run(func(ctx containerlibcontext.Context, task model.Task) error {
		var unmarshalledData map[string]interface{}
		err := json.Unmarshal(task.Attributes, &unmarshalledData)
		if err != nil {
			logrus.Error(err)
			return err
		}

		time.Sleep(15 * time.Second)
		fmt.Println(unmarshalledData)

		ctx.Success()
		ctx.RecordResult(task.Attributes)
		return nil
	})
}
