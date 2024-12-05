package parser

import "fmt"

func Parse(time []byte) string {
	fmt.Println(len(time))
	fmt.Println(time)

	for _, val := range time {
		fmt.Print(val)
	}

	return string(time)
}
