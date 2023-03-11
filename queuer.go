package queue

import (
	"log"
	"sync"
	"time"
)

type Queuer struct {
	mutex         sync.RWMutex
	tasks         []Task
	maxActiveTask int
}

func NewQueuer(maxSize int) *Queuer {
	if maxSize == 0 {
		maxSize = 10
	}
	return &Queuer{
		mutex:         sync.RWMutex{},
		tasks:         []Task{},
		maxActiveTask: maxSize,
	}
}

func (q *Queuer) Push(task Task) error {
	q.mutex.Lock()
	q.tasks = append(q.tasks, task)
	q.mutex.Unlock()
	return nil
}

func (q *Queuer) pop() {
	q.mutex.Lock()
	q.tasks = q.tasks[1:]
	q.mutex.Unlock()
}

func (q *Queuer) first() Task {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if len(q.tasks) == 0 {
		return nil
	}

	return q.tasks[0]
}

func (q *Queuer) Run() {
	activeTaskChan := make(chan Task, q.maxActiveTask)

	go func() {
		for {
			task := q.first()
			if task == nil {
				<-time.After(2 * time.Second)
				continue
			}

			activeTaskChan <- task
			log.Printf("task %s registered", task.GetTitle())
			q.pop()

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
