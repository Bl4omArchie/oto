package pkg

import (
	"fmt"
	"encoding/json"
)

func GetExecId(name string, version string) string {
	return fmt.Sprintf("%s - %s", name, version)
}

func SliceToJson(v []string) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}
