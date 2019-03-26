package workers

type worker struct {
	tasks <-chan *TaskTree
	stop  <-chan struct{}
}
type TaskTree struct {
	TaskId int64       `json:"task_id"`
	Task   []int64     `json:"task"`
	Child  []*TaskTree `json:"child""`
}
type TaskResult struct {
	TaskId int64 `json:"task_id"`
	Sum    int64 `json:"sum"`
}

func (w *worker) run() {
	for {
		select {
		case <-w.stop:
			return
		case task := <-w.tasks:
			w.doWork(task)
		}
	}
}

func initWorker() {
	for i := uint(0); i < 3; i++ {
		w := &worker{
			tasks: pool.tasks,
			stop:  pool.stop,
			pool:  pool,
		}
		go w.run()
	}
}

func (w *worker) doWork(task *TaskTree) (result TaskResult) {
	for _, item := range task.Task {
		result.Sum += item
	}
	result.TaskId = task.TaskId
	return result
}
