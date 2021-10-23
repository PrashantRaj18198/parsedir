package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {

	// testRegex()

}

func testRegex() {
	var re = regexp.MustCompile(`(?m)(\{\{\s*range\s*[^\s]*\s*\}\})`)
	var str = `{{ range .hello }}
{{ range .hello.text }}
{{         range .h.o.x.z              }}
{{range .z                    }}
{{ .name  }}`

	for i, match := range re.FindAllString(str, -1) {
		fmt.Fprintf(os.Stdout, "'%v' found at index '%d'\n", match, i)
	}
}
