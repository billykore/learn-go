package main

func main() {
	l := newTodoList()
	m := newMenu(l)
	m.start()
}
