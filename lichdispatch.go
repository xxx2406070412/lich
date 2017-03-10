// lichdispatch
package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/daryl/zeus"
)

var pstMux *zeus.Mux = nil

func init() {
	pstMux = zeus.New()
	pstMux.POST("/lich/nova", nova)
	pstMux.GET("/lich/nova", nova)
	pstMux.NotFound = notFound

}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Write(lichmsg_exception(URLUNMATCH, "url not match"))
}

//接收http 增 删 查 （改的话就是再次插入即可）
func urlParamToMap(r *http.Request) *map[string]string {

	paraMap := make(map[string]string)

	s := strings.Split(r.URL.RawQuery, "&")

	for i := 0; i < len(s); i++ {
		if "" == s[i] {
			continue
		}

		paramarray := strings.Split(s[i], "=")

		if 2 != len(paramarray) {
			continue
		}
		paraMap[paramarray[0]], _ = url.QueryUnescape(paramarray[1])
	}

	return &paraMap
}

func nova(w http.ResponseWriter, r *http.Request) {
	var act string
	act = r.URL.Query().Get("act")
	p_paraMap := urlParamToMap(r)

	var msgI msgcommon
	var pst_lichmsg = &lichmsg{p_paraMap, 0, "", "", time.Now().UnixNano() / 1e6}

	switch act {
	case "1001", "1002", "1003":
		msgI = &msgact{*pst_lichmsg, true}
	default:
		w.Write(lichmsg_exception(LOSTACT, "act code canot deal"))
		return
	}

	msgI.dealProcess()

	str, _ := json.Marshal(*(msgI.getResult()))

	w.Write(str)

}
