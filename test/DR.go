package main
import (
	"probestar/dynamicrouter"
	"bufio"
	"os"
	"log"
	"strings"
)

func main() {
	config := dynamicrouter.NewConfig("localhost:2181/gotest", "probestar", "wyw")
	r := dynamicrouter.Newdynamicrouter(config)
	model := dynamicrouter.NewRouterModel("probestar", "ps://probestar.com")
	r.Register(model)

	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		index := strings.IndexAny(command, " ");
		var c, p string
		if (index > 0) {
			c = string(data[:index])
			p = string(data[index + 1:])
		} else {
			c = command
			p = ""
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
			}()

			switch c {
			case "stop":
				running = false
			case "get":
				s := r.GetURL(p)
				if (s == "") {
					log.Println("No Router for Key: " + p)
				} else {
					log.Println(s)
				}
			case "list":
				log.Println(r.List())
			case "register":
				ps := strings.Split(p, ",")
				m := dynamicrouter.NewRouterModel(ps[0], ps[1])
				r.Register(m)
			default:
				log.Println("Unkown Command: " + c)
			}
		}()
	}
	log.Println("Done.")
}