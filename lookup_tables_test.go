package learning

import (
	"fmt"
	"testing"
)

func message(people int) {
	switch people {
	case 0:
		fmt.Println("Not a party. You are alone.")
	case 1:
		fmt.Println("One is the loneliest number")
	case 2:
		fmt.Println("Not lonely, but not a party")
	case 3:
		fmt.Println("Now we are talking")
	case 4:
		fmt.Println("Ah, yeah...")
	default:
		fmt.Println("Whoa, violated fire code!")
	}
}

func TestSwitchCase(t *testing.T) {
	message(0)
	message(1)
	message(2)
	message(3)
	message(4)
	message(5)
}

const numMessages = 6

var messages = [numMessages]string{
	"Not a party. You are alone.",
	"One is the loneliest number",
	"Not lonely, but not a party",
	"Now we are talking",
	"Ah, yeah..",
	"Whoa, violated fire code!",
}

func messageLUTs(people int) {
	if people > numMessages-1 {
		fmt.Println(messages[numMessages-1])
	} else {
		fmt.Println(messages[people])
	}
}

func TestLUTs(t *testing.T) {
	messageLUTs(0)
	messageLUTs(1)
	messageLUTs(2)
	messageLUTs(3)
	messageLUTs(4)
	messageLUTs(5)
}

func switchCase(a int) string {
	switch a {
	case 1:
		return "One"
	case 2:
		return "Two"
	case 3:
		return "Three"
	default:
		return "More than three"
	}
}

func useMap(a int) string {
	msg := map[int]string{
		1: "One",
		2: "Two",
		3: "Tree",
	}
	if s, ok := msg[a]; ok {
		return s
	}
	return "More than three"
}

const numMsgs = 3

var msgs = [numMsgs]string{
	"One",
	"Two",
	"Tree",
}

func lookupTables(a int) string {
	if a > numMsgs-1 {
		return msgs[numMsgs-1]
	}
	return msgs[a]
}

func TestLTUs(t *testing.T) {
	msg := lookupTables(2)
	t.Log(msg)
}

func BenchmarkSwitchCase(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		switchCase(i)
	}
}

func BenchmarkUseMap(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		useMap(i)
	}
}

func BenchmarkLookUpTables(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		lookupTables(i)
	}
}
