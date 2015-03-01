package queue

import(
	"container/list"
	"errors"
)

type Queue struct {
	l *list.List
}

func New() *Queue{
	return &Queue{list.New()}
}

func (q *Queue)Len() int{
	return q.l.Len()
}

func (q *Queue) Add(v interface{}) *list.Element{
	return q.l.PushBack(v)
}

func (q *Queue) Empty() bool{
	return q.l.Len() == 0
}

func (q *Queue) Head() (*list.Element, error){
	if q.Len() == 0 {
		err := errors.New("the queue is empty")
		return nil, err
	}

	e := q.l.Front()
	q.l.Remove(e)
	return e, nil
}

func (q *Queue) Clear() {
	q.l.Init()
}