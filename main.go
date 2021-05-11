package main

import "fmt"
import "github.com/json-iterator/go"
func main() {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data := 1
	marshal, err := json.Marshal(&data)
	if err != nil {
		return
	}
	fmt.Printf("json: %+v", marshal)
	fmt.Println("Hello world")
}
