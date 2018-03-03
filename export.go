package main

import (
	"fmt"
	"strings"

	"github.com/toolateforteddy/arbitrary"
	"github.com/toolateforteddy/errortrace"
)

func FormatForShellExport(data interface{}) ([]string, error) {

	flattened, err := arbitrary.FlattenWithJoiner(data, exportingJoiner)
	if err != nil {
		return nil, errortrace.Wrap(err)
	}
	retArr := make([]string, 0, len(flattened))
	for k, v := range flattened {
		retArr = append(
			retArr,
			fmt.Sprintf("%s=%s", k, v))
	}
	return retArr, nil
}

func exportingJoiner(arr []string) string {
	for i, str := range arr {
		arr[i] = strings.ToUpper(str)
	}
	return strings.Join(arr, "_")
}
