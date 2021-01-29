package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"syscall/js"

	"github.com/komkom/toml"
	"github.com/pkg/errors"
)

var Document = js.Global().Get("document")

const (
	Clear    = `input#clear`
	TOMLArea = `toml`
	JSONArea = `json`
	ErrorMsg = `errormsg`
)

func main() {

	js.Global().Set(`clear2`, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("clear\n")
		clear()
		return nil
	}))

	js.Global().Set(`format`, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("format2\n")
		load()
		return nil
	}))

	<-make(chan struct{})
}

func load() {

	errMsg := Document.Call("getElementById", ErrorMsg)
	//errMsg.Set(`innerHTML`, `testtest`)
	//style := errMsg.Get(`style`)
	//style.Set(`display`, `none`)

	j, err := transform()
	if err != nil {
		errMsg.Set(`innerHTML`, err.Error())
		return
	}

	errMsg.Set(`innerHTML`, ``)
	Document.Call("getElementById", JSONArea).Set(`innerHTML`, j)
}

func clear() {
	Document.Call("getElementById", ErrorMsg).Set(`innerHTML`, ``)
	Document.Call("getElementById", JSONArea).Set(`innerHTML`, ``)
	Document.Call("getElementById", TOMLArea).Set(`innerHTML`, ``)
}

func transform() (string, error) {

	var edit string
	val := Document.Call("getElementById", TOMLArea).Get(`value`)
	if val.Truthy() {
		edit = val.String()
	}

	r := strings.NewReader(edit)
	rd := toml.New(r)
	data, err := ioutil.ReadAll(rd)

	if err != nil {
		return ``, err
	}

	if !json.Valid(data) {
		return ``, fmt.Errorf(`generated josn not valid`)
	}

	str, err := PrettyJSON(data)
	if err != nil {
		return ``, errors.Wrap(err, `transform prettyJSON failed`)
	}

	return str, nil
}

func PrettyJSON(jsn []byte) (string, error) {

	var pretty bytes.Buffer
	err := json.Indent(&pretty, jsn, "", "&nbsp;&nbsp;&nbsp;")
	if err != nil {
		return ``, err
	}

	return string(pretty.Bytes()), nil
}
