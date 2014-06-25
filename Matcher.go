// Matcher
package Router

import (
	"net/http"
	"strings"
)

type Params struct {
	Value string
}

//Custom handler func to allow for param injection
type HandlerFunc func(*Params) ([]byte, int)

//路径及相对应的处理函数
type Matcher struct {
	//路径:   /ID/:id
	Path string
	//此模式路径相应的处理函数
	Handler HandlerFunc
}

func (m *Matcher) Process(params *Params, w http.ResponseWriter, req *http.Request) {
	content, status := m.Handler(params)
	w.WriteHeader(status)
	w.Write(content)
}

//实际URL与存储的路径进行匹配
func (m *Matcher) Matching(u string) (bool, *Params) {
	var match bool = true
	Value := ""

	is := strings.Split(u, "/")
	ps := strings.Split(m.Path, "/")

	for i := 0; i < len(is); i++ {
		if len(ps) == i || len(is) != len(ps) {
			match = false
			break
		}
		index := strings.Index(ps[i], ":")
		if is[i] != ps[i] && index != 0 {
			match = false
			break
		}
		if index == 0 {
			Value = is[i]
		}
	}

	return match, &Params{Value}
}
