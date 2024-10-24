package main

import (
	"errors"
	"fmt"
	"time"
)

type Todo struct {
	Title       string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// defining slice of Todo as Todos
type Todos []Todo

func (todos *Todos) add(title string) {
	todo := Todo{
		Title:       title,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
	*todos = append(*todos, todo)
}

func (todos *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*todos) {
		err := errors.New("invalid index")
		fmt.Println(err)
		return err
	}
	return nil
}

func (todos *Todos) delete(index int) error {
	if err := (*todos).validateIndex(index); err != nil {
		return err
	}
	*todos = append((*todos)[:index], (*todos)[index+1:]...)
	return nil
}

func (todos *Todos) toggle(index int) error {
	if err := (*todos).validateIndex(index); err != nil {
		return err
	}
	isCompleted := (*todos)[index].Completed
	if !(isCompleted) {
		(*todos)[index].Completed = true
		timeIs := time.Now()
		((*todos)[index].CompletedAt) = &timeIs
	} else {
		(*todos)[index].Completed = false
	}
	return nil
}

func (todos *Todos) edit(index int, title string) error {
	if err := (*todos).validateIndex(index); err != nil {
		return err
	}
	(*todos)[index].Title = title
	return nil
}

// func (todos *Todos) print() {
// 	table := table.New(os.Stdout)
// 	table.SetRowLines(false)
// 	table.SetHeaders("#", "Title", "Completed", "Created At", "Completed At")

// 	for index, t := range *todos {
// 		completed := emoji.CrossMark.String()
// 		completedAt := ""
// 		if t.Completed {
// 			completed = emoji.CheckMark.String()
// 			if t.CompletedAt != nil {
// 				completedAt = t.CompletedAt.Format(time.RFC1123)
// 			}
// 		}

// 		table.AddRow(strconv.Itoa(index), t.Title, completed, t.CreatedAt.Format(time.RFC1123), completedAt)
// 	}
// 	table.Render()
// }
