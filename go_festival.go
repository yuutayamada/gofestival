package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Options struct {
	Text string
	File string
	Help bool
	Loop bool
}

func main() {
	opts := cmdParse()
	if len(opts.Text) == 0 && len(opts.File) == 0 || opts.Help {
		flag.Usage()
		fmt.Println("Text or file not found.\n",
			"You can specify by -text or -file option.")
	} else {
		if len(opts.File) != 0 {
			speaker(readFile(opts.File), opts.Loop)
		} else if len(opts.Text) != 0 {
			speaker(opts.Text, opts.Loop)
		}
	}
}

func readFile(fileName string) string {
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(bs)
}

func speaker(text string, loop bool) {
	speak := func(text string) {
		ch := make(chan string, 1)
		documents := strings.Split(text, "\n")
		for len(documents) > 0 {
			documents = setter(ch, documents)
			say(ch)
		}
	}
	if loop {
		for {
			speak(text)
		}
	} else {
		speak(text)
	}
}

func setter(ch chan<- string, list []string) []string {
	ch <- list[0]
	return list[1:]
}

func cmdParse() Options {
	h := flag.Bool("h", false, "show help")
	help := flag.Bool("help", false, "show help")
	text := flag.String("text", "", "text for festival")
	loop := flag.Bool("loop", false, "whether do looping")
	file := flag.String("file", "", "file name")
	flag.Parse()
	return Options{*text, *file, *help || *h, *loop}
}

func say(document <-chan string) {
	content := strings.Replace(<-document, "\"", " ", -1)
	fmt.Println(content)
	c1 := exec.Command("echo", content)
	c2 := exec.Command("festival", "--tts")
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r
	var out bytes.Buffer
	c2.Stdout = &out
	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	c2.Wait()
	io.Copy(os.Stdout, &out)
}
