package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"syscall/js"
	"time"
)

var r, _ = regexp.Compile("(WARN|ERROR|INFO) (.+)")

type ParsedLogs = map[string]interface{}
type OutputCallback = func(parsedLogs ParsedLogs)

func ModuleOutput(parsedLogs ParsedLogs) {
	logCallback := js.Global().Get("logCallback")
	logCallback.Invoke(parsedLogs)
}

type Log struct {
	Level string
	Msg   string
}

func (l Log) ToMap() ParsedLogs {
	m := make(map[string]interface{})
	m["level"] = l.Level
	m["msg"] = l.Msg
	return m
}

func parse(str string) Log {
	log := Log{}
	if !r.MatchString(str) {
		log.Msg = str
		return log
	}
	groups := r.FindStringSubmatch(str)

	log.Level = groups[1]
	log.Msg = groups[2]
	return log
}

type Accumulator struct {
	sb  strings.Builder
	out chan string
}

func NewAccumulator() Accumulator {
	return Accumulator{sb: strings.Builder{}, out: make(chan string, 10)}
}

func (a *Accumulator) Append(str string) {
	if r.MatchString(str) {
		a.Flush()
	}
	a.sb.WriteString(str)
}

func (a *Accumulator) Flush() {
	if a.sb.Len() > 0 {
		res := a.sb.String()
		go func(res string) {
			a.out <- res
		}(res)
		a.sb = strings.Builder{}
	}
}

func Execute(callbackFn OutputCallback) {
	file, err := os.Open("./test.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	acc := NewAccumulator()
	defer close(acc.out)

	go func() {
		for {
			select {
			case str := <-acc.out:
				l := parse(str)
				callbackFn(l.ToMap())
			case <-time.After(2 * time.Second):
				acc.Flush()
			}
		}
	}()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			log.Fatal(err)
		}
		acc.Append(line)
	}

}

func main() {
	Execute(ModuleOutput)
}
