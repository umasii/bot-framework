package activityapi

import (
	"fmt"
	"sort"
	"sync"
)

var statusLogger *TaskStatusLogger

// Sets a single global status logger object to maintain instead of many as in previous implementation
// Let this be declared during runtime.
func SetStatusLogger(logger *TaskStatusLogger) {
	statusLogger = logger
}

// Gets the status logger
func GetStatusLogger() *TaskStatusLogger {
	return statusLogger
}

func CreateEmptyLogger() *TaskStatusLogger {
	return &TaskStatusLogger{
		tasksIdStatusMap: make(map[int]string),
		taskStatuses:     []TaskStatus{},
	}
}

type TaskStatus struct {
	Id     int
	Status string
	Level string
}

type internalStatus struct {
	Key   string
	Value int
}

// For: Sorting map[string]int (ideally should be moved somewhere else)
// Taken from https://stackoverflow.com/questions/18695346/how-to-sort-a-mapstringint-by-its-values
type internalStatusArr []internalStatus

func (p internalStatusArr) Len() int           { return len(p) }
func (p internalStatusArr) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p internalStatusArr) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func sortByOccurance(taskStatusToSort map[string]int) internalStatusArr {
	pl := make(internalStatusArr, len(taskStatusToSort))
	i := 0
	for k, v := range taskStatusToSort {
		pl[i] = internalStatus{k, v}
		i++
	}

	sort.Sort(sort.Reverse(pl))
	return pl
}

// Single instance checkout stats and task statuses logger.
//
// This is mainly used for CLI, but can be implemented into GRPC/UI
type TaskStatusLogger struct {
	// session checkout stats
	checkoutMutex   sync.Mutex
	CheckoutCounter int
	DeclineCounter  int
	CaptchaCounter  int
	ThreeDSCounter  int
	CartedCounter  int
	QueuePassCounter  int

	// status logging
	statusMutex      sync.Mutex
	statusMapMutex   sync.Mutex
	tasksIdStatusMap map[int]string
	taskStatuses     []TaskStatus
}

// +1 to checkout counter
func (l *TaskStatusLogger) AddCheckout() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.CheckoutCounter++
}

// +1 to decline counter
func (l *TaskStatusLogger) AddDecline() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.DeclineCounter++
}

// +1 to cart counter
func (l *TaskStatusLogger) AddCart() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.CartedCounter++
}

// -1 to cart counter
func (l *TaskStatusLogger) RemoveCart() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.CartedCounter--
}

// +1 to queue pass counter
func (l *TaskStatusLogger) AddQueuePass() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.QueuePassCounter++
}

// -1 to queue pass counter
func (l *TaskStatusLogger) RemoveQueuePass() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.QueuePassCounter--
}

// +1 to captcha counter
func (l *TaskStatusLogger) AddCaptcha() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.CaptchaCounter++
}

// +1 to decline counter
func (l *TaskStatusLogger) Add3DS() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.ThreeDSCounter++
}

// -1 to captcha counter
func (l *TaskStatusLogger) RemoveCaptcha() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.CaptchaCounter--
}

// -1 to decline counter
func (l *TaskStatusLogger) Remove3DS() {
	l.checkoutMutex.Lock()
	defer l.checkoutMutex.Unlock()

	l.ThreeDSCounter--
}

// Returns task status array and empties the array
func (l *TaskStatusLogger) GetStatuses() []TaskStatus {
	l.statusMutex.Lock()
	defer l.statusMutex.Unlock()

	copiedTaskStatuses := l.taskStatuses
	l.taskStatuses = []TaskStatus{}
	return copiedTaskStatuses
}

// Pushes new task status update to array
func (l *TaskStatusLogger) PushTaskStatusUpdate(id int, status string, level string) {
	l.statusMutex.Lock()
	defer l.statusMutex.Unlock()
	l.taskStatuses = append(l.taskStatuses, TaskStatus{
		Id:     id,
		Status: status,
		Level: level,
	})
}

// Gets task status array and returns it as formatted string
func (l *TaskStatusLogger) GetTaskStatusArray() []string {
	l.statusMapMutex.Lock()
	defer l.statusMapMutex.Unlock()

	tasksMap := map[string]int{}

	statuses := l.GetStatuses()
	// Map task statuses by ID in case of multiple status updates
	for i := 0; i < len(statuses); i++ {
		l.tasksIdStatusMap[statuses[i].Id] = statuses[i].Status
	}

	// Count occurance of each status message
	for _, status := range l.tasksIdStatusMap {
		if val, ok := tasksMap[status]; ok {
			tasksMap[status] = val + 1
		} else {
			tasksMap[status] = 1
		}
	}

	// Sort by occurance
	sorted := sortByOccurance(tasksMap)

	arr := []string{}
	for _, val := range sorted {
		arr = append(arr, fmt.Sprintf("[%d] %s", val.Value, val.Key))
	}

	return arr
}
