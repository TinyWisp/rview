package main

import (
	"encoding/json"
	"fmt"

	"github.com/TinyWisp/rview/template"
)

func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

func Test() {
	// str := "111 + 222+333 * 5555 * 6666+(33*77+22) > abc && (var2 > func1(abc,def,333) || var2 > 333 || t2 || t3 == true || t4 == false || !t5) && t1 != nil"
	// str := "abc || t1 != nil"
	str := "!abc || t1 != nil"
	exps, err := template.ReadTplExp(str)
	fmt.Println(str)
	fmt.Println(err)
	PrettyPrint(exps)
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("==================================================")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	exps2, err2 := template.ParseTplExp(str)
	fmt.Println(err2)
	PrettyPrint(exps2)
}

func main() {
	Test()
}
