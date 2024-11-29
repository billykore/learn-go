package main

import "fmt"

type todoList struct {
	todos []*todo
}

func newTodoList() *todoList {
	return &todoList{}
}

func (l *todoList) getTodos() {
	fmt.Println("List of todos")
	if len(l.todos) == 0 {
		fmt.Println("(empty)")
	}
	for _, t := range l.todos {
		t.print()
	}
}

func (l *todoList) addTodos(todos ...*todo) {
	l.todos = append(l.todos, todos...)
}

func (l *todoList) removeTodo(id int) {
	for i := range l.todos {
		if i == id {
			l.todos = append(l.todos[:i], l.todos[i+1:]...)
		}
	}
}

func (l *todoList) completeTodo(id int) {
	for i, t := range l.todos {
		if i == id {
			t.setCompleted(true)
		}
	}
}

func (l *todoList) reset() {
	l.todos = []*todo{}
}
