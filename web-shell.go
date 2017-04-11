package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	//	"unicode/utf8"
	"flag"
	"time"
)

// kballard: go-shellquote

var (
	UnterminatedSingleQuoteError = errors.New("Unterminated single-quoted string")
	UnterminatedDoubleQuoteError = errors.New("Unterminated double-quoted string")
	UnterminatedEscapeError      = errors.New("Unterminated backslash-escape")
)

var (
	flagport          = flag.Int("p", 8080, "port to serve")
	splitChars        = " \n\t"
	singleChar        = '\''
	doubleChar        = '"'
	escapeChar        = '\\'
	doubleEscapeChars = "$`\"\n\\"
)

// func Split(input string) (words []string, err error) {
// 	var buf bytes.Buffer
// 	words = make([]string, 0)
//
// 	for len(input) > 0 {
// 		// skip any splitChars at the start
// 		c, l := utf8.DecodeRuneInString(input)
// 		if strings.ContainsRune(splitChars, c) {
// 			input = input[l:]
// 			continue
// 		}
//
// 		var word string
// 		word, input, err = splitWord(input, &buf)
// 		if err != nil {
// 			return
// 		}
// 		words = append(words, word)
// 	}
// 	return
// }
//
// func splitWord(input string, buf *bytes.Buffer) (word string, remainder string, err error) {
// 	buf.Reset()
//
// raw:
// 	{
// 		cur := input
// 		for len(cur) > 0 {
// 			c, l := utf8.DecodeRuneInString(cur)
// 			cur = cur[l:]
// 			if c == singleChar {
// 				buf.WriteString(input[0 : len(input)-len(cur)-l])
// 				input = cur
// 				goto single
// 			} else if c == doubleChar {
// 				buf.WriteString(input[0 : len(input)-len(cur)-l])
// 				input = cur
// 				goto double
// 			} else if c == escapeChar {
// 				buf.WriteString(input[0 : len(input)-len(cur)-l])
// 				input = cur
// 				goto escape
// 			} else if strings.ContainsRune(splitChars, c) {
// 				buf.WriteString(input[0 : len(input)-len(cur)-l])
// 				return buf.String(), cur, nil
// 			}
// 		}
// 		if len(input) > 0 {
// 			buf.WriteString(input)
// 			input = ""
// 		}
// 		goto done
// 	}
//
// escape:
// 	{
// 		if len(input) == 0 {
// 			return "", "", UnterminatedEscapeError
// 		}
// 		c, l := utf8.DecodeRuneInString(input)
// 		if c == '\n' {
// 		} else {
// 			buf.WriteString(input[:l])
// 		}
// 		input = input[l:]
// 	}
// 	goto raw
//
// single:
// 	{
// 		i := strings.IndexRune(input, singleChar)
// 		if i == -1 {
// 			return "", "", UnterminatedSingleQuoteError
// 		}
// 		buf.WriteString(input[0:i])
// 		input = input[i+1:]
// 		goto raw
// 	}
//
// double:
// 	{
// 		cur := input
// 		for len(cur) > 0 {
// 			c, l := utf8.DecodeRuneInString(cur)
// 			cur = cur[l:]
// 			if c == doubleChar {
// 				buf.WriteString(input[0 : len(input)-len(cur)-l])
// 				input = cur
// 				goto raw
// 			} else if c == escapeChar {
// 				c2, l2 := utf8.DecodeRuneInString(cur)
// 				cur = cur[l2:]
// 				if strings.ContainsRune(doubleEscapeChars, c2) {
// 					buf.WriteString(input[0 : len(input)-len(cur)-l-l2])
// 					if c2 == '\n' {
// 					} else {
// 						buf.WriteRune(c2)
// 					}
// 					input = cur
// 				}
// 			}
// 		}
// 		return "", "", UnterminatedDoubleQuoteError
// 	}
//
// done:
// 	return buf.String(), input, nil
// }
//
// // go-shellquote ends here.

func cmdExec(s string) string {
	if s == "" {
		return "need command"
	}

	// expand environmental variables
	s = os.ExpandEnv(s)

	log.Println("exec:", s)
	args := strings.Split(s, " ")

	//args, _ := Split(s)
	var cmd *exec.Cmd
	if len(args) < 2 {
		cmd = exec.Command(args[0])
	} else {
		cmd = exec.Command(args[0], args[1:]...)
	}
	var out = new(bytes.Buffer)

	cmd.Stdout = out
	cmd.Stderr = out

	err := cmd.Start()
	if err != nil {
		outputString := out.String()
		log.Println(outputString)
		log.Println("error:", err.Error())
		return htmlformat(outputString) + " " + htmlformat(err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		outputString := out.String()
		log.Println(outputString)
		log.Println("error:", err.Error())
		return htmlformat(outputString) + " " + htmlformat(err.Error())
	}
	outputString := out.String()
	log.Println(outputString)
	return htmlformat(outputString)

}

func htmlformat(s string) string {
	return strings.Replace(s, "\n", "<br />", -1)
}
func init() {
	log.SetFlags(log.Ltime)
}

func handler(w http.ResponseWriter, r *http.Request) {
	headhtml := `<!doctype html>
<head><meta charset="utf-8" />
	<title>` + os.Args[0] + `</title>
	<style>
		body {
		overflow-x: hidden;
		background-color: black;
		color: green;
		text-align: center;
		}
		.code {
		text-align: left;
		background-color: black;
		color: green;
	}

	#c {
	width: 94vw;
}
		#output {
		text-align: left;
		overflow-y: scroll;
		overflow-x: hidden;
		width: 100vw;
		height: 80vh;
	}
	</style>
</head>
<body>`

	fmt.Fprintf(w, headhtml)

	fmt.Fprintf(w, `<div id="output" class="code">`)
	if r.Method == "POST" {
		ubuf.WriteString(htmlformat("\n"))
		cmd := r.FormValue("c")
		cmd = os.ExpandEnv(cmd)
		switch {
		case cmd == "":
		case cmd == "quit":
			os.Exit(0)
		case cmd == "clear", cmd == "reset":
			ubuf.Truncate(0)
		case cmd == "cd":
			err := os.Chdir(os.Getenv("HOME"))
			if err != nil {
				ubuf.WriteString(htmlformat(err.Error() + "\n"))
			}
		case strings.HasPrefix(cmd, "cd "):
			err := os.Chdir(strings.TrimPrefix(cmd, "cd "))
			if err != nil {
				ubuf.WriteString(htmlformat(err.Error() + "\n"))
			}
		default:
			log.Println(r.RemoteAddr, "command:", cmd)
			ubuf.WriteString(cmdExec(cmd) + "\n")
		}
	}
	fmt.Fprintf(w, ubuf.String())
	fmt.Fprintf(w, `</div>`)
	formhtml := `` +
		`	<div class="code"><form method="post"><input type="text" name="c" id="c" class="code"/>
		<br />
		<input type="submit" value="Ok" />
	</form></div>
	<script>
	var objDiv = document.getElementById("output");
	objDiv.scrollTop = objDiv.scrollHeight;


	</script>
	<script>document.getElementById('c').focus();</script>
	</body>
</html>`
	fmt.Fprintf(w, formhtml)
}

var ubuf = new(bytes.Buffer)

func main() {
	flag.Parse()
	http.HandleFunc("/", handler)

	go func() {
		<-time.After(time.Second)
		println("Listening on port:", *flagport)
	}()
	err := http.ListenAndServe(fmt.Sprint(":", *flagport), nil)
	if err != nil {
		println(err)
		os.Exit(111)
	}
}
