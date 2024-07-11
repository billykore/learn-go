package main

import "testing"

func TestPrintTodo(t *testing.T) {
	td := newTodo("title", "description")
	td.print()
}

func TestTodoList_AddTodo(t *testing.T) {
	td1 := newTodo("title", "description")
	td2 := newTodo("title", "description")
	td3 := newTodo("title", "description")
	list := newTodoList()
	list.addTodos(td1, td2, td3)
	if len(list.todos) != 3 {
		t.Error("expected 3 todos")
	}
}

func TestTodoList_RemoveTodo(t *testing.T) {
	td1 := newTodo("title", "description")
	td2 := newTodo("title", "description")
	td3 := newTodo("title", "description")
	list := newTodoList()
	list.addTodos(td1, td2, td3)

	list.removeTodo(0)
	if len(list.todos) != 2 {
		t.Error("expected 2 todos")
	}

	list.getTodos()
}

func TestTodoList_GetTodos(t *testing.T) {
	td1 := newTodo("title", "description")
	td2 := newTodo("title", "description")
	td3 := newTodo("title", "description")
	list := newTodoList()
	list.addTodos(td1, td2, td3)
	list.getTodos()
}

func TestMenu_Display(t *testing.T) {
	l := newTodoList()
	m := newMenu(l)
	m.display()
}
