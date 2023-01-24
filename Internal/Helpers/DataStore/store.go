package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/kardianos/osext"
	"github.com/cicadaaio/LVBot/Internal/Errors"
)

// Data should be in the same directory as the built.exe
func GetStoreFilePath(storeName string) string {
	b, _ := osext.Executable()
	storePath := path.Join(path.Dir(b), "Data", fmt.Sprintf("%s.json", storeName))
	return storePath
}

func Write(data interface{}, storeName string) {
	file, _ := json.MarshalIndent(data, "", " ")

	err := ioutil.WriteFile(GetStoreFilePath(storeName), file, 0644)

	if err != nil {
		Errors.Handler(err)
	}
}

func Read(data interface{}, storeName string) {
	file, err := ioutil.ReadFile(GetStoreFilePath(storeName))

	if err != nil {
		Errors.Handler(err)
	}

	json.Unmarshal([]byte(file), data)
}
