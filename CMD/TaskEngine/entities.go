package taskengine

import (
	"sync"

	tasks "github.com/umasii/bot-framework/internal/tasks"
)

type TaskEngine struct {
	TaskGroups []BotTaskGroup
}

type BotTaskGroup struct {
	GroupID       int
	GroupName     string
	GroupStatus   string
	QtyTasks      int
	QtyActive     int
	QtyCheckedOut int
	QtyFailed     int
	Tasks         []tasks.IBotTask
	TaskWaitGroup sync.WaitGroup
}

func (btg *BotTaskGroup) StartTasks(wg *sync.WaitGroup) {
	for i := range btg.Tasks {
		wg.Add(1)
		go btg.Tasks[i].WrapExecutor(btg.Tasks[i].Execute, wg)
	}
	wg.Wait()
	//TODO: Make sure we call wg.Done() whenever the task is done or stopped!
}

func (btg *BotTaskGroup) StopTasks() {
	for i := range btg.Tasks {
		go btg.Tasks[i].Stop()
	}
}