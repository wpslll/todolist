package todo

import (
	"backend/psql"
	"errors"
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
	task.Id = taskDB.Id
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

func (l *List) AddTask(task Task) (Task, error) {
	taskDto, err := l.db.Insert(task.Title, task.Description, task.IsCompleted, task.CreatedAt, task.CompletedAt)
	if err != nil {
		return Task{}, err
	}
	var tmp Task
	tmp = convertToTask(taskDto)
	return tmp, nil
}

func (l *List) CompleteTask(id int) (Task, error) {
	time := time.Now()
	v, err := l.db.Complete_Uncomplete(id, true, &time)
	if err != nil {
		return Task{}, err
	}
	task := convertToTask(v)
	return task, nil
}

func (l *List) UncompleteTask(id int) (Task, error) {
	v, err := l.db.Complete_Uncomplete(id, false, nil)
	if err != nil {
		return Task{}, err
	}
	task := convertToTask(v)
	return task, nil
}

func (l *List) DeleteTask(id int) error {
	if err := l.db.Delete(id); err != nil {
		return err
	}
	return nil
}

func (l *List) NewUser(login, password string) error {
	if err := l.db.InsertUser(login, password); err != nil {
		return err
	}
	return nil
}

func (l *List) FindUser(login, password string) (int, error) {
	id, err := l.db.SelectUser(login, password)
	if err != nil {
		return -1, err
	}
	if id == -1 {
		return -1, errors.New("No such user")
	}
	return id, nil
}
func (l *List) FindUserId(id int) error {
	return l.db.SelectUserId(id)
}