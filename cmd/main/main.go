package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/iveronanomi/toad/config"
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
	flag.StringVar(&formatFlag, "f", "//TODO({{username}}).{{project}} {{issue_type}}", "TODO format string")
	flag.Parse()

	log.SetFlags(log.Lshortfile)
	if !verboseFlag {
		log.SetOutput(&bytes.Buffer{})
	}

	if pathFlag == "" {
		flag.Usage()
		os.Exit(1)
	}
}

const chunkSize = 1024 //* 1024 * 1024

func main() {
	defer config.Save()

	excludes = strings.Split(excludesFlag, ",")
	paths = strings.Split(pathFlag, ",")
	extensions = strings.Split(typesFlag, ",	")

	log.Printf("Paths %#v", paths)
	log.Printf("Excluded %#v", excludes)
	log.Printf("Extensions %#v", extensions)

	cr := tracker.New(tracker.GitCreator)
	cr.SetAccessToken("")

	_, err := config.Read()
	if err != nil {
		panic(err)
	}

	chunks := make(chan []string)
	col := todo.New(paths, excludes, extensions, chunkSize)
	go col.Collect(chunks)
	go read(chunks)
	<-quit

	for _, l := range list {
		for _, f := range l {
			f.Close()
		}
	}

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

var list map[string]map[int]*os.File

func find(paths []string) {
	for _, path := range paths {

		f, err := os.Open(path)
		defer func() {
			if err != nil {
				f.Close()
			}
		}()

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
				list[path][line] = f
			}
		}

		if err := scanner.Err(); err != nil {
			// Handle the error
		}
	}
}

func create() {

}
