package queue

type Task interface {
	GetTitle() string
	Do() error
}
