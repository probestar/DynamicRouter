package dynamicrouter
import "bytes"

type RouterModel struct {
	Key string
	URL string
}

func NewRouterModel(key string, url string) *RouterModel {
	return &RouterModel{Key:key, URL:url}
}

func (model *RouterModel)String() string {
	var buf bytes.Buffer
	buf.WriteString("Key: ")
	buf.WriteString(model.Key)
	buf.WriteString("; URL: ")
	buf.WriteString(model.URL)
	return buf.String()
}