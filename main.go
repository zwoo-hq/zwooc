package main

import (
	"fmt"

	"github.com/zwoo-hq/zwoo-builder/pkg/config"
)

func main() {
	c, err := config.Load("./zwoo.config.json")
	fmt.Println(err)
	// fmt.Println(c)

	for _, v := range c.Profiles {
		fmt.Print("Profile: ")
		fmt.Print(v.Name())
		fmt.Println(v.BaseOptions())
	}
}
