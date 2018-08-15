package main

/**
    Author: luzequan
    Created: 2018-01-02 11:21:18
*/
import (
	"drcs/exec/web"

	//_ "drcs/dep/nodelib/dep"
	//_ "drcs/dep/nodelib/batchdistribution"
	_ "drcs/dep/nodelib/crp"

	//_ "net/http/pprof"
	//"net/http"
	//"fmt"
	//"io"
)

func main() {


	//go func() {
	//	http.HandleFunc("/", index)
	//	if err := http.ListenAndServe("0.0.0.0:9090", nil); err != nil {
	//	}
	//}()


	web.Run()
}

//func index(w http.ResponseWriter, _ *http.Request) {
//	w.Header().Set("Content-type", "text/html")
//	io.WriteString(w, "<h2>Links</h2>\n<ul>")
//	for _, link := range []string{"/advance", "/simple"} {
//		fmt.Fprintf(w, `<li><a href="%v">%v</a>`, link, link)
//	}
//	io.WriteString(w, "</ul>")
//}

