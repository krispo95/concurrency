package workers

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

type worker struct {
	tasks   chan *TaskTree
	results chan *TaskResult
}

func newWorker(tasks chan *TaskTree, results chan *TaskResult) *worker {
	return &worker{tasks: tasks, results: results}
}

type TaskTree struct {
	TaskId uint        `json:"task_id"`
	Task   []uint      `json:"task"`
	Child  []*TaskTree `json:"child""`
}
type TaskResult struct {
	TaskId uint `json:"task_id"`
	Sum    uint `json:"sum"`
}

//var Tasks = make(chan TaskTree, 10)
//var Result = make(chan TaskResult, 10)

func readJSON(r io.Reader) *TaskTree {
	tree := new(TaskTree)
	err := json.NewDecoder(r).Decode(tree)
	if err != nil {
		panic("Ошибка преобразования")
	}
	return tree
}

func (w *worker) run(parRes uint) {
	for task := range w.tasks {
		res := w.doWork(task)
		res.Sum += parRes
		w.results <- &res
		if task.Child == nil {
			return
		}

		for _, item := range task.Child {
			w.tasks <- item
		}
		return
	}
}

func initWorker(tasks chan *TaskTree, results chan *TaskResult, workerCount uint) {
	for i := uint(0); i < workerCount; i++ {
		w := newWorker(tasks, results)
		go w.run(0)
	}
}

func (w *worker) doWork(task *TaskTree) (result TaskResult) {
	for _, item := range task.Task {
		result.Sum += item
	}
	time.Sleep(time.Second)
	result.TaskId = task.TaskId
	return result
}

func InitWork(path string) {
	time0 := time.Now()
	file, err := os.Open(path)
	if err != nil {
		panic("Ошибка чтения файла")
	}
	tree := readJSON(file)
	var tasks = make(chan *TaskTree)
	var results = make(chan *TaskResult)
	initWorker(tasks, results, 3)
	tasks <- tree
	readResult(results)
	log.Printf("duration: %+v", time.Now().Sub(time0))
}

func readResult(results chan *TaskResult) {
	for {
		log.Printf("%+v", <-results)
	}
}

//func doWork(task *TaskTree) (result TaskResult) {
//	for _, item := range task.Task {
//		result.Sum += item
//	}
//	time.Sleep(time.Second)
//	result.TaskId = task.TaskId
//	return result
//}
//func concurrWork(tasks *TaskTree, parRes int64) {
//
//	res := doWork(tasks)
//	res.Sum += parRes
//	Result <- res
//
//	if tasks.Child == nil {
//		return
//	}
//
//	for _, item := range tasks.Child {
//		go concurrWork(item, res.Sum)
//	}
//	return
//}
//func InitWork(path string) {
//	time0 := time.Now()
//	file, err := os.Open(path)
//	if err != nil {
//		panic("Ошибка чтения файла")
//	}
//	tree := readJSON(file)
//	concurrWork(tree, 0)
//	result := []TaskResult{}
//	for item := range Result {
//		result = append(result, item)
//	}
//	log.Printf("result: %+v", result)
//	log.Printf("duration: %+v", time.Now().Sub(time0))
//	log.Printf("Result: %+v", Result)
//}

//-----------------------------------------------
//func doWork(task *TaskTree) (result TaskResult) {
//	for _, item := range task.Task {
//		result.Sum += item
//	}
//	time.Sleep(time.Second)
//	result.TaskId = task.TaskId
//	return result
//}
//
//func seqWork(tasks *TaskTree, parRes int64) (result []TaskResult) {
//
//	res := doWork(tasks)
//	res.Sum += parRes
//	result = append(result, res)
//
//	if tasks.Child == nil {
//		return result
//	}
//
//	for _, item := range tasks.Child {
//		childRes := seqWork(item, res.Sum)
//		result = append(result, childRes...)
//	}
//	return result
//}
//
//func InitWork(path string) {
//	time0 := time.Now()
//	file, err := os.Open(path)
//	if err != nil {
//		panic("Ошибка чтения файла")
//	}
//	tree := readJSON(file)
//	result := seqWork(tree, 0)
//	log.Printf("duration: %+v", time.Now().Sub(time0))
//	log.Printf("Result: %+v", result)
//}
//------------------------------------------------
