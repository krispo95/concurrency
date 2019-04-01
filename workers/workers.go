package workers

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

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

var Tasks = make(chan TaskTree, 10)
var Result = make(chan TaskResult, 10)

func readJSON(r io.Reader) *TaskTree {
	tree := new(TaskTree)
	err := json.NewDecoder(r).Decode(tree)
	if err != nil {
		panic("Ошибка преобразования")
	}
	return tree
}

func doWork(task *TaskTree) (result TaskResult) {
	for _, item := range task.Task {
		result.Sum += item
	}
	time.Sleep(time.Second)
	result.TaskId = task.TaskId
	return result
}
func concurrWork(tasks *TaskTree, parRes int64) {

	res := doWork(tasks)
	res.Sum += parRes
	Result <- res

	if tasks.Child == nil {
		return
	}

	for _, item := range tasks.Child {
		go concurrWork(item, res.Sum)
	}
	return
}
func InitWork(path string) {
	time0 := time.Now()
	file, err := os.Open(path)
	if err != nil {
		panic("Ошибка чтения файла")
	}
	tree := readJSON(file)
	concurrWork(tree, 0)
	result := []TaskResult{}
	for item := range Result {
		result = append(result, item)
	}
	log.Printf("result: %+v", result)
	log.Printf("duration: %+v", time.Now().Sub(time0))
	log.Printf("Result: %+v", Result)
}

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
//func (w *worker) run() {
//	for {
//		select {
//		case <-w.stop:
//			return
//		case task := <-w.tasks:
//			w.doWork(task)
//		}
//	}
//}
//
//func initWorker() {
//	for i := uint(0); i < 3; i++ {
//		w := &worker{
//			tasks: pool.tasks,
//			stop:  pool.stop,
//			pool:  pool,
//		}
//		go w.run()
//	}
//}
//
//func (w *worker) doWork(task *TaskTree) (result TaskResult) {
//	for _, item := range task.Task {
//		result.Sum += item
//	}
//	result.TaskId = task.TaskId
//	return result
//}
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
