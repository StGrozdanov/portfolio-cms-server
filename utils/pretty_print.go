package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(input interface{}) {
	bytes, _ := json.MarshalIndent(input, "", "")
	fmt.Println(string(bytes))
}
