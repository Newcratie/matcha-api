package logprint

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"strings"
)

const pad = 35
const format = "---%*s%*s---\n"

func PrettyPrint(v interface{}) (err error) {
	pretty := color.New(color.FgCyan).Add(color.BgHiBlack)
	b, err := json.MarshalIndent(v, "", "  ")
	temp := strings.Split(string(b), "\n")
	if err == nil {
		for _, line := range temp {
			pretty.Printf("%s\n", line)
		}
	}
	return
}

func Title(s string) {
	title := color.New(color.FgBlack).Add(color.BgHiBlue)
	title.Printf(format, pad+len(s)/2, s, pad-len(s)/2, "")
}

func Centered(s string) {
	fmt.Printf(format, pad+len(s)/2, s, pad-len(s)/2, "")
}

func Error(err error) {
	s := "Error: " + err.Error()
	title := color.New(color.FgHiRed)
	title.Printf("%*s%*s---\n", 1, s, pad*2-len(s)+3, "")
}

func End() {
	s := "endlog"
	end := color.New(color.FgBlack).Add(color.BgHiBlack)
	end.Printf(format, pad+len(s)/2, s, pad-len(s)/2, "")
}
