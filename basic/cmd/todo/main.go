package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dyxj/grpc-exp/basic/todo"
	"github.com/golang/protobuf/proto"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: list or add")
		os.Exit(1)
	}

	var err error
	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list()
	case "add":
		err = add(strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const (
	dbPath = "mydb.pb"
)

func add(text string) error {
	task := &todo.Task{
		Text: text,
		Done: false,
	}

	tl := todo.TaskList{}
	tl.Tlist = append(tl.Tlist, task)

	b, err := proto.Marshal(&tl)
	if err != nil {
		return fmt.Errorf("could not encode task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", dbPath, err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close file %s: %v", dbPath, err)
	}
	return nil
}

func list() error {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	var tasklist todo.TaskList
	if err := proto.Unmarshal(b, &tasklist); err != nil {
		return fmt.Errorf("could not read tasklist: %v", err)
	}

	for _, v := range tasklist.Tlist {
		if v.Done {
			fmt.Printf("Done!")
		} else {
			fmt.Printf("Not Done!")
		}
		fmt.Printf(" %s\n", v.Text)
	}

	return nil
}
