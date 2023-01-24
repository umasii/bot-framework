package TaskStore

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/cicadaaio/LVBot/Internal/Errors"
	store "github.com/cicadaaio/LVBot/Internal/Helpers/DataStore"
	"github.com/cicadaaio/LVBot/Internal/Tasks"
)

var taskGroupList []Tasks.TaskGroup

func AddTask(task Tasks.IBotTask) error {
	taskGroup, err := GetTaskGroupByID(task.Get().GroupID)
	if err != nil {
		return Errors.Handler(errors.New("Failed to Add Task"))
	}

	setNewTaskID(taskGroup, &task)

	taskGroup.Tasks = append(taskGroup.Tasks, task)
	store.Write(&taskGroupList, "tasks")
	return nil
}

func AddTaskGroup(taskGroup *Tasks.TaskGroup) {
	taskGroups := GetTaskGroups()
	taskGroup.GroupID = getNextTaskGroupID(&taskGroups)
	if taskGroup.Tasks == nil {
		taskGroup.Tasks = []Tasks.IBotTask{}
	}
	taskGroups = append(taskGroups, *taskGroup)
	store.Write(&taskGroups, "tasks")
}

func GetTaskByID(groupID, taskID int) (*Tasks.IBotTask, error) {
	taskGroup, err := GetTaskGroupByID(groupID)

	if err != nil {
		return nil, Errors.Handler(errors.New(fmt.Sprintf("Could not locate Task Group (GroupID: %v)", groupID)))
	}

	for i := range taskGroup.Tasks {
		if taskGroup.Tasks[i].Get().TaskID == taskID {
			return &taskGroup.Tasks[i], nil
		}
	}

	return nil, Errors.Handler(errors.New(fmt.Sprintf("Could not locate Task (GroupID: %v, TaskID: %v)", groupID, taskID)))

}

func GetTaskGroups() []Tasks.TaskGroup {
	taskGroups, err := CustomUnmarshalTasksFile()

	if err != nil {
		Errors.Handler(err)
	}

	return taskGroups
}

func GetTaskGroupByID(taskGroupID int) (*Tasks.TaskGroup, error) {
	taskGroupList = GetTaskGroups()

	for i := range taskGroupList {
		if taskGroupList[i].GroupID == taskGroupID {
			return &taskGroupList[i], nil
		}
	}

	return nil, Errors.Handler(errors.New(fmt.Sprintf("Task Group with ID %v was not found", taskGroupID)))
}

func setNewTaskID(taskGroup *Tasks.TaskGroup, task *Tasks.IBotTask) {
	taskID := reflect.ValueOf(*task).Elem()
	field := taskID.FieldByName("TaskID")
	if field.IsValid() {
		field.SetInt(int64(getNextTaskID(&taskGroup.Tasks)))
	}
}

func getNextTaskID(tasks *[]Tasks.IBotTask) int {
	if len(*tasks) > 0 {
		lastTask := (*tasks)[len(*tasks)-1]
		return lastTask.Get().TaskID + 1
	}
	return 1
}

func getNextTaskGroupID(taskGroups *[]Tasks.TaskGroup) int {
	if len(*taskGroups) > 0 {
		lastTaskGroup := (*taskGroups)[len(*taskGroups)-1]
		return lastTaskGroup.GroupID + 1
	} else {
		return 1
	}
}
