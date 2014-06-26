// Processor
package Router

import (
	"net/http"
)

type Params struct {
	Value string
}

//Custom handler func to allow for param injection
type HandlerFunc func(*Params) ([]byte, int)

type Processor struct {
	params *Params
	//此模式路径相应的处理函数
	Handler HandlerFunc
}

func (p *Processor) Process(w http.ResponseWriter, req *http.Request) {
	content, status := p.Handler(p.params)
	w.WriteHeader(status)
	w.Write(content)
}
