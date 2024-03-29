package main

import (
	todo "Todo"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var todoFileName = ".todo.json"

func main() {

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool, Developed for the pragmatic bookshelf. And im copying it to learn go \n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2049 wakaka\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information : ")
		flag.PrintDefaults()

	}
	// Parsing Command Line Flag
	add := flag.Bool("add", false, "Add Task to The Todo List")
	list := flag.Bool("list", false, "List all task")
	complete := flag.Int("complete", 0, "Item To be Completed")

	flag.Parse()

	l := &todo.List{}

	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch {
	case *list:
		// List current to do items
		fmt.Print(l)
	case *complete > 0:
		// Complete the given item
		if err := l.Completed(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		// When any arguments (excluding flags) are provided, they will be
		// used as the new task
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		l.Add(t)
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, ""), nil
	}
	s := bufio.NewScanner(r)
	s.Scan()

	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}
	return s.Text(), nil
}
