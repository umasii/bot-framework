package utilities

import "encoding/json"

// Pretty prints a struct
//
// Better than fmt.Sprintf("%#v", i)
func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s) + "\n"
}
