package TaskEngine

import (
	"github.com/cicadaaio/LVBot/CMD/DataStores/TaskStore"
	"github.com/cicadaaio/LVBot/Internal/Tasks"
)

func (te *TaskEngine) loadTasksIntoEngine(taskGroups *[]Tasks.TaskGroup) {
	var botTaskGroups = []BotTaskGroup{}

	for tgIndex := range *taskGroups {

		currentTaskGroup := BotTaskGroup{
			GroupName:   (*taskGroups)[tgIndex].GroupName,
			GroupID:     (*taskGroups)[tgIndex].GroupID,
			QtyTasks:    len((*taskGroups)[tgIndex].Tasks),
			GroupStatus: "idle",
			Tasks:       (*taskGroups)[tgIndex].Tasks,
		}

		botTaskGroups = append(botTaskGroups, currentTaskGroup)
	}

	te.TaskGroups = botTaskGroups
}

func createTasksFromJSON() *[]Tasks.TaskGroup {
	taskGroupsData := TaskStore.GetTaskGroups()
	taskGroups := []Tasks.TaskGroup{}

	for i := range taskGroupsData {

		currentTaskGroup := Tasks.TaskGroup{
			GroupName: taskGroupsData[i].GroupName,
			GroupID:   taskGroupsData[i].GroupID,
			Tasks:     taskGroupsData[i].Tasks,
		}

		taskGroups = append(taskGroups, currentTaskGroup)
	}

	return &taskGroups
}
