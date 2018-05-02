package main

import "github.com/k0kubun/pp"

type data struct {
	name map[string]string
}

func getName(d data) {
	pp.Println(d.name)
	d.name["ddd"] = "derek"

}

func main() {
	da := new(data)
	da.name = make(map[string]string, 0)
	da.name["ddd"] = "test"
	getName(*da)
	pp.Println(da.name)
}
