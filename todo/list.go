package todo

import (
	"backend/psql"
	"time"
)

type List struct {
	db psql.DataBase
}

func NewList(Db psql.DataBase) *List {
	return &List{
		db: Db,
	}
}

func convertToTask(taskDB psql.TaskDto) Task {
	var task Task
	task.Title = taskDB.Title
	task.Description = taskDB.Description
	task.IsCompleted = taskDB.IsCompleted
	task.CreatedAt = taskDB.CreatedAt
	task.CompletedAt = taskDB.CompletedAt
	return task
}

func (l *List) ListTasks() (map[string]Task, error) {
    dbTasks, err := l.db.Select()
    if err != nil {
        return nil, err 
    }
    
    tmp := make(map[string]Task)
    for k, v := range dbTasks {
        tmp[k] = convertToTask(v)
	}
    return tmp, nil
}

func (l *List) GetTask(title string) (Task, error) {
	v, err := l.db.Select_title(title)
	if err != nil {
		return Task{}, err
	}
	task := convertToTask(v)
	return task, nil
}

func (l *List) ListUncompletedTasks() (map[string]Task, error) {
	tmp := make(map[string]Task)
	dbTasks, err := l.db.Select_uncompleted()
	if err != nil {
		return nil, err
	}
	for k, v := range dbTasks {
		tmp[k] = convertToTask(v)
	}
	return tmp, nil
}

func (l *List) AddTask(task Task) error {
	err := l.db.Insert(task.Title, task.Description, task.IsCompleted, task.CreatedAt, task.CompletedAt)
	if err != nil {
		return err
	}
	return nil
}

func (l *List) CompleteTask(title string) (Task, error) {
	if err := l.db.Complete_Uncomplete(title, true, time.Now()); err != nil {
		return Task{}, err
	}
	v, err := l.db.Select_title(title)
	if err != nil {
		return Task{}, err
	}
	task := convertToTask(v)
	return task, nil
}

func (l *List) UncompleteTask(title string) (Task, error) {
	if err := l.db.Complete_Uncomplete(title, false, time.Now()); err != nil {
		return Task{}, err
	}
	v, err := l.db.Select_title(title)
	if err != nil {
		return Task{}, err
	}
	task := convertToTask(v)
	return task, nil
}

func (l *List) DeleteTask(title string) error {
	if err := l.db.Delete(title); err != nil {
		return err
	}
	return nil
}
