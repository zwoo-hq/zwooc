package main

import (
	"fmt"

	"github.com/zwoo-hq/zwoo-builder/pkg/config"
)

type (
	Foo struct {
		Bar string `json:"bar"`
	}
)

func main() {
	c, err := config.Load("./zwoo.config.json")
	fmt.Println(err)
	// fmt.Println(c)
	profiles, _ := c.GetProfiles()
	for _, v := range profiles {
		fmt.Print("Profile: ")
		fmt.Print(v.Name())
		fmt.Println(v)
	}

	fragments, _ := c.GetFragments()
	for _, v := range fragments {
		fmt.Print("Fragment: ")
		fmt.Print(v.Name())
		fmt.Println(v)
	}

	compounds, _ := c.GetCompounds()
	for _, v := range compounds {
		fmt.Print("Compound: ")
		fmt.Print(v.Name())
		fmt.Println(v)
	}
}
