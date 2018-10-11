package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var (
	watcher  *fsnotify.Watcher
	basePath string
)

func init() {
	w, err := fsnotify.NewWatcher()

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create the watcher: %s\n", err)
		os.Exit(1)
	}

	watcher = w

	if args := strings.Join(os.Args, " "); strings.Contains(args, "-h") || strings.Contains(args, "--help") {
		usage()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "missing path as argument\n")
		os.Exit(1)
	}

	basePath = filepath.Join(filepath.Dir(os.Args[0]), os.Args[len(os.Args)-1])

	if err := lookForProtoc(); err != nil {
		fmt.Fprintf(os.Stderr, "could not find the protoc library, please ensure you have it installed: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	registerListeners(basePath)
	watch()
}

func lookForProtoc() error {
	cmd := exec.Command("protoc", "--version")
	return cmd.Start()
}

func usage() {
	out, err := exec.Command("protoc", "--help").Output()
	if err != nil {
		fmt.Fprint(os.Stderr, "could not read: protoc --help\n")
		os.Exit(1)
	}
	fmt.Print(strings.Replace(string(out), "protoc ", "protoc-watch ", 1))
}

func registerListeners(path string) {
	watcher.Add(path)
	descriptors, err := ioutil.ReadDir(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read path: %s\n", err)
		os.Exit(1)
	}

	for _, descriptor := range descriptors {
		next := filepath.Join(path, descriptor.Name())

		if descriptor.IsDir() && string(descriptor.Name()[0]) != "." {
			registerListeners(next)
			continue
		}

		watcher.Add(next)
	}
}

func watch() {
	for {
		select {
		case event := <-watcher.Events:
			handle(event)
		case err := <-watcher.Errors:
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func handle(event fsnotify.Event) {
	if event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Chmod == fsnotify.Chmod {
		return
	}

	if event.Op&fsnotify.Remove == fsnotify.Remove {
		watcher.Remove(event.Name)
		return
	}

	descriptor, err := os.Stat(event.Name)

	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read path: %s\n", err)
	}

	if descriptor.IsDir() {
		registerListeners(event.Name)
		return
	}

	if !strings.Contains(event.String(), ".proto") {
		return
	}

	compile(event.Name)
}

func compile(path string) {
	args := strings.Replace(strings.Join(os.Args[1:], " "), basePath, path, 2)
	args = strings.Replace(args, ".proto/", ".proto", 1)
	out, _ := exec.Command("protoc", strings.Split(args, " ")...).Output()
	fmt.Println(path, string(out))
}
