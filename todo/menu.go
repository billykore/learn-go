package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type menu struct {
	todoList  *todoList
	keepGoing bool
}

func newMenu(list *todoList) *menu {
	return &menu{
		todoList:  list,
		keepGoing: true,
	}
}

const displayText = `
Welcome to Todo List App
Select an option:
  1. Display Todo List
  2. Add Todo
  3. Delete Todo
  4. Complete Todo
  5. Reset
  0. Exit
`

func (m *menu) display() {
	fmt.Print(displayText)
}

func (m *menu) start() {
	for m.keepGoing {
		m.display()
		m.processOption()
	}
}

func (m *menu) processOption() {
	option, err := m.getInput("Select an option: ")
	if err != nil {
		log.Fatalf("error scanning option: %v", err)
	}
	m.useOption(option)
}

func (m *menu) useOption(option string) {
	switch option {
	case "1":
		m.displayTodoListOption()
		break
	case "2":
		m.addTodoOption()
		break
	case "3":
		m.deleteTodoOption()
		break
	case "4":
		m.completeTodoOption()
		break
	case "5":
		m.resetOption()
		break
	case "0":
		m.exitOption()
		break
	}
}

func (m *menu) displayTodoListOption() {
	m.todoList.getTodos()
}

func (m *menu) addTodoOption() {
	title, err := m.getInput("Title: ")
	if err != nil {
		log.Fatalf("error getting Title: %v", err)
	}
	desc, err := m.getInput("Description: ")
	if err != nil {
		log.Fatalf("error getting Title: %v", err)
	}
	t := newTodo(title, desc)
	m.todoList.addTodos(t)
}

func (m *menu) deleteTodoOption() {
	strId, err := m.getInput("ID: ")
	if err != nil {
		log.Fatalf("error scanning id: %v", err)
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Fatalf("error convert strId to id: %v", err)
	}
	m.todoList.removeTodo(id)
}

func (m *menu) completeTodoOption() {
	strId, err := m.getInput("ID: ")
	if err != nil {
		log.Fatalf("error scanning id: %v", err)
	}
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Fatalf("error convert strId to id: %v", err)
	}
	m.todoList.completeTodo(id)
}

func (m *menu) getInput(title string) (string, error) {
	fmt.Print(title)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		return "", fmt.Errorf("error reading input %s: %v", title, err)
	}
	return scanner.Text(), nil
}

func (m *menu) resetOption() {
	m.todoList.reset()
}

func (m *menu) exitOption() {
	m.keepGoing = false
}
