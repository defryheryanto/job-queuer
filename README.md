# job-queuer
A Simple Job Queuer for Go

### Create a new queuer
Below code will return a queuer with 10 max worker
```
maxWorker := 10
jobQueuer := queuer.NewQueuer(maxWorker)
```

### Running the queue
Refer to this code
```
jobQueuer.Run(context.Background())
```
This will run the queue in a separate goroutine (in background), any Task that will be pushed before or after the Run() function will be executed.
This function WILL NOT locking

### Adding Task to queue
This package will use the Task interface for the queued object

For the first step, let's make a struct that implements the Task interface<br>
`GetTitle()` will be logged into the terminal when the Task is queued, being run, and done (failed and finished)<br>
`Do(context.Context)` will be executed when the Task is being run<br>
```
type SimpleTask struct {}
func(s *SimpleTask) GetTitle() string {
  return "simple_task"
}
func (s *SimpleTask) Do(ctx context.Context) error {
	time.Sleep(1 * time.Second)
	return nil
}
```

The second step, we create a SimpleTask object and pass it into the queuer
```
simpleTask := &SimpleTask{}
jobQueuer.Push(simpleTask)
```
This code will push a SimpleTask (implements the Task interface) to the queuer. The task will be executed according to the queue and when there is an available worker.
