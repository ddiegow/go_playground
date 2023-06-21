package capitalize

import (
	"fmt"
	"strings"
)

func Capitalize(input []string) {
	for i := range input {
		fmt.Println(strings.ToUpper(input[i]))
	}
}
