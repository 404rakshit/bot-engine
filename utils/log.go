package utils

import (
	"encoding/json"
	"fmt"
)

func Print(v ...any) {
	const (
		green = "\033[92m"
		reset = "\033[0m"
	)

	for _, val := range v {
		b, err := json.MarshalIndent(val, "", "  ")
		if err != nil {
			fmt.Printf("%s%v%s\n", green, val, reset)
			continue
		}
		fmt.Printf("%s%s%s\n", green, b, reset)
	}
}
