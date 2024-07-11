package main

import (
	"fmt"
)

type todo struct {
	title       string
	description string
	completed   bool
}

func newTodo(title, description string) *todo {
	return &todo{
		title:       title,
		description: description,
		completed:   false,
	}
}

const printFormat = `
Title:       %s
Description: %s
Completed:   %v
`

func (t *todo) print() {
	fmt.Printf(printFormat, t.title, t.description, t.completed)
}

func (t *todo) setCompleted(completed bool) {
	t.completed = completed
}
