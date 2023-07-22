package taskengine

import (
	"github.com/umasii/bot-framework/cmd/datastores/taskstore"
	"github.com/umasii/bot-framework/internal/tasks"
)

func (te *TaskEngine) loadTasksIntoEngine(taskGroups *[]tasks.TaskGroup) {
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

func createTasksFromJSON() *[]tasks.TaskGroup {
	taskGroupsData := taskstore.GetTaskGroups()
	taskGroups := []tasks.TaskGroup{}

	for i := range taskGroupsData {

		currentTaskGroup := tasks.TaskGroup{
			GroupName: taskGroupsData[i].GroupName,
			GroupID:   taskGroupsData[i].GroupID,
			Tasks:     taskGroupsData[i].Tasks,
		}

		taskGroups = append(taskGroups, currentTaskGroup)
	}

	return &taskGroups
}
