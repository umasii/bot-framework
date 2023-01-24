package TaskEngine

import (
	"errors"
	"fmt"
	"sync"

	"github.com/cicadaaio/LVBot/Internal/Tasks"
)




func (te *TaskEngine) InitializeEngine() {
	createdTaskGroups := createTasksFromJSON()
	te.loadTasksIntoEngine(createdTaskGroups)
}

func (te *TaskEngine) StartAllTasks() {
	for i := range te.TaskGroups {
		te.TaskGroups[i].StartTasks(&te.TaskGroups[i].TaskWaitGroup)
	}
}

func (te *TaskEngine) StopAllTasks() {
	for i := range te.TaskGroups {
		te.TaskGroups[i].StopTasks()
	}
}

func (te *TaskEngine) GetRunningTaskCount() int {
	taskCount := 0
	for i := range te.TaskGroups {
		taskCount += len(te.TaskGroups[i].Tasks)
	}

	return taskCount
}

func (te *TaskEngine) StartSelectedTasks(wg *sync.WaitGroup, tasks *[]Tasks.IBotTask) {
	for i := range *tasks {
		wg.Add(1)
		go (*tasks)[i].WrapExecutor((*tasks)[i].Execute, wg)
	}
	wg.Wait()
}

func (te *TaskEngine) StartTasksInGroup(groupID int) {
	for i := range te.TaskGroups {
		if te.TaskGroups[i].GroupID == groupID {
			te.TaskGroups[i].StartTasks(&te.TaskGroups[i].TaskWaitGroup)
		}
	}
}

func (te *TaskEngine) GetTasksInTaskGroup(groupID int) ([]Tasks.IBotTask, error){
	for i := range te.TaskGroups {
		if te.TaskGroups[i].GroupID == groupID {
			return te.TaskGroups[i].Tasks, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Task Group with ID %v not found!", groupID))
}