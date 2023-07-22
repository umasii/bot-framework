package taskstore

import (
	goErrors "errors"
	"fmt"
	"reflect"

	errors "github.com/umasii/bot-framework/internal/errors"
	store "github.com/umasii/bot-framework/internal/helpers/datastore"
	tasks "github.com/umasii/bot-framework/internal/tasks"
)

var taskGroupList []tasks.TaskGroup

func AddTask(task tasks.IBotTask) error {
	taskGroup, err := GetTaskGroupByID(task.Get().GroupID)
	if err != nil {
		return errors.Handler(goErrors.New("Failed to Add Task"))
	}

	setNewTaskID(taskGroup, &task)

	taskGroup.Tasks = append(taskGroup.Tasks, task)
	store.Write(&taskGroupList, "tasks")
	return nil
}

func AddTaskGroup(taskGroup *tasks.TaskGroup) {
	taskGroups := GetTaskGroups()
	taskGroup.GroupID = getNextTaskGroupID(&taskGroups)
	if taskGroup.Tasks == nil {
		taskGroup.Tasks = []tasks.IBotTask{}
	}
	taskGroups = append(taskGroups, *taskGroup)
	store.Write(&taskGroups, "tasks")
}

func GetTaskByID(groupID, taskID int) (*tasks.IBotTask, error) {
	taskGroup, err := GetTaskGroupByID(groupID)

	if err != nil {
		return nil, errors.Handler(goErrors.New(fmt.Sprintf("Could not locate Task Group (GroupID: %v)", groupID)))
	}

	for i := range taskGroup.Tasks {
		if taskGroup.Tasks[i].Get().TaskID == taskID {
			return &taskGroup.Tasks[i], nil
		}
	}

	return nil, errors.Handler(goErrors.New(fmt.Sprintf("Could not locate Task (GroupID: %v, TaskID: %v)", groupID, taskID)))

}

func GetTaskGroups() []tasks.TaskGroup {
	taskGroups, err := CustomUnmarshalTasksFile()

	if err != nil {
		errors.Handler(err)
	}

	return taskGroups
}

func GetTaskGroupByID(taskGroupID int) (*tasks.TaskGroup, error) {
	taskGroupList = GetTaskGroups()

	for i := range taskGroupList {
		if taskGroupList[i].GroupID == taskGroupID {
			return &taskGroupList[i], nil
		}
	}

	return nil, errors.Handler(goErrors.New(fmt.Sprintf("Task Group with ID %v was not found", taskGroupID)))
}

func setNewTaskID(taskGroup *tasks.TaskGroup, task *tasks.IBotTask) {
	taskID := reflect.ValueOf(*task).Elem()
	field := taskID.FieldByName("TaskID")
	if field.IsValid() {
		field.SetInt(int64(getNextTaskID(&taskGroup.Tasks)))
	}
}

func getNextTaskID(tasks *[]tasks.IBotTask) int {
	if len(*tasks) > 0 {
		lastTask := (*tasks)[len(*tasks)-1]
		return lastTask.Get().TaskID + 1
	}
	return 1
}

func getNextTaskGroupID(taskGroups *[]tasks.TaskGroup) int {
	if len(*taskGroups) > 0 {
		lastTaskGroup := (*taskGroups)[len(*taskGroups)-1]
		return lastTaskGroup.GroupID + 1
	} else {
		return 1
	}
}
