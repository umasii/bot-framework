package TaskStore

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/cicadaaio/LVBot/Internal/Constants"
	"github.com/cicadaaio/LVBot/Internal/Errors"
	store "github.com/cicadaaio/LVBot/Internal/Helpers/DataStore"
	"github.com/cicadaaio/LVBot/Internal/Tasks"
)

type TaskGroupUnpacker struct {
	GroupName string
	GroupID   int
	Tasks     []TaskUnpacker
}

type TaskUnpacker struct {
	Site  string
	Value Tasks.IBotTask
}

func (c *TaskUnpacker) UnmarshalJSON(data []byte) error {

	value, err := UnmarshalCustomValue(data, "Site", Constants.SITES)

	if err != nil {
		return Errors.Handler(err)
	}

	c.Value = value

	return nil
}

func UnmarshalCustomValue(data []byte, typeJsonField string, customTypes map[string]reflect.Type) (Tasks.IBotTask, error) {
	m := map[string]interface{}{}

	if err := json.Unmarshal(data, &m); err != nil {
		return nil, Errors.Handler(err)
	}

	typeName := m[typeJsonField].(string)

	var value Tasks.IBotTask
	if ty, found := customTypes[typeName]; found {
		value = reflect.New(ty).Interface().(Tasks.IBotTask)
	}

	valueBytes, err := json.Marshal(m)

	if err != nil {
		return nil, Errors.Handler(err)
	}

	if err = json.Unmarshal(valueBytes, &value); err != nil {
		return nil, Errors.Handler(err)
	}

	return value, Errors.Handler(err)
}

func readTasksJSONFile() ([]byte, error) {
	jsonFile, err := os.Open(store.GetStoreFilePath("tasks"))

	if err != nil {
		return []byte{}, Errors.Handler(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue, nil
}

func CustomUnmarshalTasksFile() ([]Tasks.TaskGroup, error) {
	taskGroups := []Tasks.TaskGroup{}

	jsonData, err := readTasksJSONFile()

	if err != nil {
		return taskGroups, Errors.Handler(err)
	}

	var unpackedTaskGroups []TaskGroupUnpacker

	err = json.Unmarshal(jsonData, &unpackedTaskGroups)

	if err != nil {
		return taskGroups, Errors.Handler(err)
	}

	for i := range unpackedTaskGroups {

		currentTaskGroup := Tasks.TaskGroup{
			GroupName: unpackedTaskGroups[i].GroupName,
			GroupID:   unpackedTaskGroups[i].GroupID,
			Tasks:     []Tasks.IBotTask{},
		}

		for j := range unpackedTaskGroups[i].Tasks {
			currentTaskGroup.Tasks = append(currentTaskGroup.Tasks, unpackedTaskGroups[i].Tasks[j].Value)
		}

		taskGroups = append(taskGroups, currentTaskGroup)

	}

	return taskGroups, nil
}
