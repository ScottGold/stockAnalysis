package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	help   bool
	infile string
)

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.StringVar(&infile, "i", "", "input file")
}

func usage() {
	fmt.Fprintf(os.Stderr,
		`app version: app/1.10.0 Usage:  [-h] 
Options:
`)
	flag.PrintDefaults()
}

func readline(filename string, cline chan string) error {
	ifile, err := os.Open(filename)
	if err != nil {
		return err
	}

	go func() {
		line := ""
		bytes := make([]byte, 102400)
		count, err := ifile.Read(bytes)
		for count > 0 {
			fmt.Println("read", count, err)
			lastlinestart := 0
			for i := 0; i < count; i++ {
				if bytes[i] == '\n' {
					line = line + string(bytes[lastlinestart:i])
					cline <- line
					line = ""
					lastlinestart = i + 1
				}
				if i == count-1 && i > lastlinestart {
					line = string(bytes[lastlinestart:i])
				}
			}
			count, err = ifile.Read(bytes)
		}
		if line != "" {
			cline <- line
		}
		close(cline)
	}()

	return nil
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	chRLine := make(chan string, 100)
	err := readline(infile, chRLine)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	finished := make(chan bool)
	go func() {
		l, ok := <-chRLine
		for ; ok == true; l, ok = <-chRLine {
			fmt.Println(l)

		}
		finished <- true
	}()
	<-finished
}
