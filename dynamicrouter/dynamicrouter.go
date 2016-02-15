package dynamicrouter
import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
	"log"
	"encoding/json"
	"strings"
	"bytes"
)

type Router interface {
	List() string
	Register(*RouterModel)
	GetURL(string) string
}

type dynamicRouter struct {
	config *Config
	conn   *zk.Conn
	models map[string][]RouterModel
}

func Newdynamicrouter(config *Config) Router {
	conn, event := connectServer(config)
	r := &dynamicRouter{config: config, conn:conn, models:map[string][]RouterModel{}}
	r.handleEvent(event)
	r.refresh()
	return r
}

func (router *dynamicRouter) Register(model *RouterModel) {
	key := router.config.Path() + "/" + model.Key
	buf, _ := json.Marshal(model)
	s, _ := router.conn.Create(key, buf, 3, []zk.ACL{{zk.PermAll, "world", "anyone"}})
	log.Println("Register as " + s)
}

func (router *dynamicRouter) GetURL(key string) string {
	count := len(router.models[key])
	if (count == 0) {
		return ""
	}
	index := (int)(time.Now().UTC().UnixNano()) % count
	return router.models[key][index].URL
}

func (router *dynamicRouter) List() string {
	var buf bytes.Buffer
	for k, v := range router.models {
		buf.WriteString(k)
		buf.WriteString("{")
		for i, r := range v {
			buf.WriteString(r.URL)
			if (i != len(v) - 1) {
				buf.WriteString(", ")
			}
		}
		buf.WriteString("}\n")
	}
	return buf.String()
}

func connectServer(config *Config) (*zk.Conn, zk.Event) {
	conn, e, _ := zk.Connect([]string{config.Address()}, time.Second)
	event := <-e
	return conn, event
}

func (router *dynamicRouter) handleEvent(event zk.Event) {
	log.Printf("Got Event: %+v", event)
	switch event.Type {
	case zk.EventNodeChildrenChanged:
		router.refresh()
	}
}

func (router *dynamicRouter) handleDeleteEvent(event zk.Event) {
	index := strings.LastIndexAny(event.Path, "/")
	key := string([]byte(event.Path)[index + 1:])
	log.Printf("Node [%s] is offline.", key)
}

func (router *dynamicRouter) refresh() {
	log.Print("Refresh " + router.config.Path())
	router.models = map[string][]RouterModel{}
	go func() {
		children, _, e, _ := router.conn.ChildrenW(router.config.Path())
		for _, value := range children {
			go func(value string) {
				buf, _, _ := router.conn.Get(router.config.Path() + "/" + value)
				var model RouterModel
				json.Unmarshal(buf, &model)
				s := router.models[model.Key]
				if (s == nil) {
					s = []RouterModel{}
				}
				router.models[model.Key] = append(s, model)
				log.Printf("Load from " + value + ". " + model.String())
			}(value)
		}
		event := <-e
		router.handleEvent(event)
	}()
}