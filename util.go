package ckRequest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/rodaine/table"
)

func jsonErr(err error, data []byte) {
	if err != nil {
		ioutil.WriteFile(kJsonParseErrFilePath, data, 0666)
		printParseErrDetailMsg(data, nil)
		panicIf(err)
	}
}

func reactIf(err error, f func()) {
	if err != nil {
		f()
		panicIf(err)
	}
}

func jsonToMap(data []byte) map[string]interface{} {
	var d map[string]interface{}
	jsonErr(json.Unmarshal(data, &d), data)
	return d
}

func jsonToStruct(data []byte, v interface{}) {
	jsonErr(json.Unmarshal(data, v), data)
}

func jsonConvert(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func isHtml(data []byte) bool {
	return bytes.Contains(data, []byte("<!DOCTYPE html")) ||
		bytes.Contains(data, []byte("<html>")) ||
		bytes.Contains(data, []byte("<body>")) ||
		bytes.Contains(data, []byte("<div>"))
}

func isJson(data []byte) bool {
	return data[0] == byte('{')
}

type jsonParseErrCallbackFunc func(data []byte) (details string, suggest string)

// func defaultCb(data []byte) (details string, suggest string) { return }

func printParseErrDetailMsg(data []byte, cb jsonParseErrCallbackFunc) {
	color.Red("Error detected: json parse")
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tab := table.New("Data type", "Error details", "Suggest").WithWidthFunc(runewidth.StringWidth)
	tab.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	dataType := "unknown"
	details := "null"
	suggest := "检查" + kJsonParseErrFilePath
	if cb != nil {
		if isHtml(data) {
			dataType = "html"
			details, suggest = cb(data)
		} else if isJson(data) {
			dataType = "json"
			details, suggest = cb(data)
		}
	}
	tab.AddRow(dataType, details, suggest)
	tab.Print()
}

func panicIf(err error) {
	if err == nil {
		return
	}

	color.Red("Error catched: %s", err.Error())
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tab := table.New("# Traceback", "Func", "# Line")
	tab.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for i := 0; i < 10; i++ {
		pc, _, line, _ := runtime.Caller(i)
		f := runtime.FuncForPC(pc)
		tab.AddRow("Error stack"+itoa(i), f.Name(), line)
		if f.Name() == "main.main" || f.Name() == "runtime.goexit" {
			break
		}
	}

	tab.Print()
	os.Exit(1)
}

func itoa(number int) string {
	return strconv.Itoa(number)
}

func atoi(number string) (int, error) {
	return strconv.Atoi(number)
}

func atoiMust(number string) int {
	i, err := strconv.Atoi(number)
	panicIf(err)
	return i
}
