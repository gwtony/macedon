package test
import (
	"io"
	//"os"
	//"fmt"
	"net/http"
	//"io/ioutil"
	//"path/filepath"
	"net/http/httptest"
	//"git.lianjia.com/lianjia-sysop/napi/log"
)
type TLog struct {
}
func (l *TLog) Info(arg0 interface{}, args ...interface{}) {}
func (l *TLog) Debug(arg0 interface{}, args ...interface{}) {}
func (l *TLog) Error(arg0 interface{}, args ...interface{}) {}

type Thandler struct {
}
func (t *Thandler)ServeHTTP (w http.ResponseWriter, req *http.Request) {
}

const Test_napi_conf string = `
[default]
addr: 172.30.23.39:8888

log: napi.log
level: debug

location: /dns
`

func TestInitlog() (* TLog) {
	return &TLog{}
}

func Test_generate_rr(method, uri string, body io.Reader) (*httptest.ResponseRecorder, *http.Request){
	r, _ := http.NewRequest(method, uri, body)
	w := httptest.NewRecorder()
	return w, r
}
