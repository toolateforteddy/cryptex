package export

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/toolateforteddy/arbitrary"
	"github.com/toolateforteddy/errortrace"
)

func LoadFile(filename string) (interface{}, error) {
	var config interface{}
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		fmt.Println(err)
		return nil, errortrace.Wrap(err)
	}
	return config, nil
}

func FormatForShellExport(data interface{}) ([]string, error) {

	flattened, err := arbitrary.FlattenWithJoiner(data, exportingJoiner)
	if err != nil {
		return nil, errortrace.Wrap(err)
	}
	retArr := make([]string, 0, len(flattened))
	for k, v := range flattened {
		retArr = append(
			retArr,
			fmt.Sprintf(`%s="%v"`, k, v))
	}
	return retArr, nil
}

func exportingJoiner(arr []string) string {
	for i, str := range arr {
		arr[i] = strings.ToUpper(str)
	}
	return strings.Join(arr, "_")
}
