package main

import (
	"9fans.net/go/acme"
	_ "embed"
	"encoding/json"
	"fmt"
	"golang.org/x/term"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func logFatalIfErr[T any](value T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return value
}

//go:embed README.md
var usage string

func window(myJSON any, path string) {
	myWindow := logFatalIfErr(acme.New())

	myBody := func() (result string) {
		switch myJSON := myJSON.(type) {
		case bool, nil, float64, string:
			myWindow.Write("tag", []byte(func() string {
				switch myJSON.(type) {
				case bool:
					return "Boolean"
				case nil:
					return "Null"
				case float64:
					return "Number"
				case string:
					return "String"
				default:
					panic(nil)
				}
			}()))
			result += fmt.Sprintf("%v\n", myJSON)
		case map[string]any:
			myWindow.Write("tag", []byte("Object "))
			for key := range myJSON {
				result += key + "\n"
			}
		case []any:
			myWindow.Write("tag", []byte("Array "))
			for index := range myJSON {
				result += strconv.Itoa(index) + "\n"
			}
		}
		return
	}()

	myWindow.Name(path)
	myWindow.Write("body", []byte(myBody))
	myWindow.Ctl("clean")

	myWindow.Write("addr", []byte("0"))
	myWindow.Ctl("dot=addr")
	myWindow.Ctl("show")

	for myEvent := range myWindow.EventChan() {
		switch myEvent.C2 {
		case 'x', 'X':
			if string(myEvent.Text) == "Del" {
				myWindow.Ctl("delete")
			}
		case 'l':
			// To-do: go up?
		case 'L':
			switch myJSON := myJSON.(type) {
			case bool, nil, float64, string:
			case map[string]any, []any:
				text := string(myEvent.Text)

				value, ok := func() (any, bool) {
					switch myJSON := myJSON.(type) {
					case map[string]any:
						value, ok := myJSON[text]
						return value, ok
					case []any:
						number, err := strconv.Atoi(text)
						ok := err == nil &&
							number >= 0 &&
							number < len(myJSON)
						if !ok {
							return nil, false
						}
						return myJSON[number], true
					default:
						panic(nil)
					}
				}()

				if !ok {
					continue
				}

				waitGroup.Add(1)
				go window(
					value,
					path+"/"+text,
				)
			}
		}
	}

	myWindow.CloseFiles()
	waitGroup.Done()
}

var waitGroup sync.WaitGroup

func die() { fmt.Fprintf(os.Stderr, usage); os.Exit(1) }

func main() {
	if len(os.Args) > 2 {
		die()
	}

	var myJSON = func() any {
		jsonReader := func() io.Reader {
			// Prioritize command line arguments over standard input.
			if len(os.Args) > 1 {
				arg := string(os.Args[1])
				if len(arg) > 0 &&
					arg[0] == '-' &&
					(len(arg) == 1 ||
						(arg[1] != '0' &&
							arg[1] != '1' &&
							arg[1] != '2' &&
							arg[1] != '3' &&
							arg[1] != '4' &&
							arg[1] != '5' &&
							arg[1] != '6' &&
							arg[1] != '7' &&
							arg[1] != '8' &&
							arg[1] != '9')) {
					die()
				}
				return strings.NewReader(os.Args[1])
			}
			if term.IsTerminal(int(os.Stdin.Fd())) {
				die()
			}
			return os.Stdin
		}()

		myJSON, err := io.ReadAll(jsonReader)
		ok := err == nil && len(myJSON) != 0
		if !ok {
			die()
		}

		var result any
		json.Unmarshal(myJSON, &result)
		return result
	}()

	waitGroup.Add(1)
	go window(myJSON, "/json")
	waitGroup.Wait()
}
