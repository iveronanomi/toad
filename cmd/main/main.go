package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/iveronanomi/todo"
	"github.com/iveronanomi/todo/tracker"
)

var (
	pathFlag, excludesFlag, typesFlag, formatFlag string
	excludes, paths, extensions                   []string

	verboseFlag bool

	quit  = make(chan int, 2)
	found int
)

func init() {
	flag.StringVar(&pathFlag, "p", ".", "Source files `directories`\ntodo -p /home/projects/fancy-go-project\n")
	flag.StringVar(&excludesFlag, "e", "", "Exclude `directories`; separate with [,]")
	flag.StringVar(&typesFlag, "t", "", "Searching files `extensions`, separate with [,]")
	flag.BoolVar(&verboseFlag, "v", false, "Verbose mode")
	flag.StringVar(&formatFlag, "f", "//TODO({{username}});{{project}};{{issue_type}}", "TODO format string")
	flag.Parse()

	if !verboseFlag {
		log.SetOutput(&bytes.Buffer{})
	}

	if pathFlag == "" {
		flag.Usage()
		os.Exit(1)
	}
}

type fundamental struct {
	todo string
	line uint64
	src  string
	body string
}

var todoList []*fundamental

func todoFromSource(src string, line uint64, format string) *fundamental {
	return &fundamental{
		src:  src,
		line: line,
		body: fmt.Sprintf("issue auto-generated from todo `%s`", src),
	}
}

const chunkSize = 1024 //* 1024 * 1024

func main() {
	excludes = strings.Split(excludesFlag, ",")
	paths = strings.Split(pathFlag, ",")
	extensions = strings.Split(typesFlag, ",	")

	log.Printf("Paths %#v", paths)
	log.Printf("Excluded %#v", excludes)
	log.Printf("Extensions %#v", extensions)

	cr := tracker.New(tracker.GitCreator)
	cr.SetAccessToken("")

	url, err := tracker.Create("hello world", "body issue", "iveronanomi", "todo", cr)
	if err != nil {
		panic(err)
	}
	log.Print(url)

	os.Exit(1)

	chunks := make(chan []string)
	col := todo.New(paths, excludes, extensions, chunkSize)
	go col.Collect(chunks)
	go read(chunks)
	<-quit
	log.Print("Find: ", found)
}

func read(ch chan []string) {
	for {
		select {
		case v, ok := <-ch:
			{
				if ok {
					if len(v) < 1 {
						quit <- 0
					}
					go find(v)
					continue
				}
			}
		}
	}
}

func find(paths []string) {
	for _, path := range paths {

		f, err := os.Open(path)
		defer f.Close()

		if err != nil {
			panic(err)
		}

		// Splits on newlines by default.
		scanner := bufio.NewScanner(f)

		// https://golang.org/pkg/bufio/#Scanner.Scan
		var line int
		for scanner.Scan() {
			line++
			if strings.Contains(scanner.Text(), "//TODO") {
				todoList = append(todoList, todoFromSource(path, uint64(line), ""))
			}
		}

		if err := scanner.Err(); err != nil {
			// Handle the error
		}
	}
}

func create() {

}
