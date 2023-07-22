package taskstore

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"

	constants "github.com/umasii/bot-framework/internal/constants"
	errors "github.com/umasii/bot-framework/internal/errors"
	store "github.com/umasii/bot-framework/internal/helpers/datastore"
	tasks "github.com/umasii/bot-framework/internal/tasks"
)

type TaskGroupUnpacker struct {
	GroupName string
	GroupID   int
	Tasks     []TaskUnpacker
}

type TaskUnpacker struct {
	Site  string
	Value tasks.IBotTask
}

func (c *TaskUnpacker) UnmarshalJSON(data []byte) error {

	value, err := UnmarshalCustomValue(data, "Site", constants.SITES)

	if err != nil {
		return errors.Handler(err)
	}

	c.Value = value

	return nil
}

func UnmarshalCustomValue(data []byte, typeJsonField string, customTypes map[string]reflect.Type) (tasks.IBotTask, error) {
	m := map[string]interface{}{}

	if err := json.Unmarshal(data, &m); err != nil {
		return nil, errors.Handler(err)
	}

	typeName := m[typeJsonField].(string)

	var value tasks.IBotTask
	if ty, found := customTypes[typeName]; found {
		value = reflect.New(ty).Interface().(tasks.IBotTask)
	}

	valueBytes, err := json.Marshal(m)

	if err != nil {
		return nil, errors.Handler(err)
	}

	if err = json.Unmarshal(valueBytes, &value); err != nil {
		return nil, errors.Handler(err)
	}

	return value, errors.Handler(err)
}

func readTasksJSONFile() ([]byte, error) {
	jsonFile, err := os.Open(store.GetStoreFilePath("tasks"))

	if err != nil {
		return []byte{}, errors.Handler(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue, nil
}

func CustomUnmarshalTasksFile() ([]tasks.TaskGroup, error) {
	taskGroups := []tasks.TaskGroup{}

	jsonData, err := readTasksJSONFile()

	if err != nil {
		return taskGroups, errors.Handler(err)
	}

	var unpackedTaskGroups []TaskGroupUnpacker

	err = json.Unmarshal(jsonData, &unpackedTaskGroups)

	if err != nil {
		return taskGroups, errors.Handler(err)
	}

	for i := range unpackedTaskGroups {

		currentTaskGroup := tasks.TaskGroup{
			GroupName: unpackedTaskGroups[i].GroupName,
			GroupID:   unpackedTaskGroups[i].GroupID,
			Tasks:     []tasks.IBotTask{},
		}

		for j := range unpackedTaskGroups[i].Tasks {
			currentTaskGroup.Tasks = append(currentTaskGroup.Tasks, unpackedTaskGroups[i].Tasks[j].Value)
		}

		taskGroups = append(taskGroups, currentTaskGroup)

	}

	return taskGroups, nil
}
