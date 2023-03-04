package queue

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Queuer struct {
	mutex         sync.Mutex
	tasks         []Task
	maxActiveTask int
}

func NewQueuer(maxSize int) *Queuer {
	if maxSize == 0 {
		maxSize = 10
	}
	return &Queuer{
		mutex:         sync.Mutex{},
		tasks:         []Task{},
		maxActiveTask: maxSize,
	}
}

func (q *Queuer) Add(task Task) error {
	q.mutex.Lock()
	_, existingTask := q.getTask(task)
	if existingTask != nil {
		return fmt.Errorf("task with title %s is already exists", task.GetTitle())
	}

	q.tasks = append(q.tasks, task)
	q.mutex.Unlock()
	return nil
}

func (q *Queuer) Remove(task Task) {
	q.mutex.Lock()
	indexTask, _ := q.getTask(task)
	if indexTask == -1 {
		return
	}

	q.tasks = append(q.tasks[:indexTask], q.tasks[indexTask+1:]...)
	q.mutex.Unlock()
}

func (q *Queuer) getTask(task Task) (int, Task) {
	for i, t := range q.tasks {
		if t.GetTitle() == task.GetTitle() {
			return i, t
		}
	}

	return -1, nil
}

func (q *Queuer) getFirstTask() Task {
	if len(q.tasks) == 0 {
		return nil
	}

	return q.tasks[0]
}

func (q *Queuer) Run() {
	activeTaskChan := make(chan Task, q.maxActiveTask)

	go func() {
		for {
			task := q.getFirstTask()
			if task == nil {
				time.Sleep(2 * time.Second)
				continue
			}

			activeTaskChan <- task
			log.Printf("task %s registered", task.GetTitle())
			q.Remove(task)

			go func() {
				defer func() {
					<-activeTaskChan
				}()
				log.Printf("running task %s", task.GetTitle())
				err := task.Do()
				if err != nil {
					log.Printf("error running task %s with error %v\n", task.GetTitle(), err)
					return
				}
				log.Printf("success: %s", task.GetTitle())
			}()
		}
	}()
}
