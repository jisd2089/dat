package web
//
//import (
//	"net/http"
//	"text/template"
//	"encoding/json"
//
//	"dat/core"
//	"dat/common/session"
//	"dat/config"
//	"github.com/henrylee2cn/pholcus/logs"
//	"dat/runtime/status"
//)
//
//var globalSessions *session.Manager
//
//func init() {
//	config := `{"cookieName":"pholcusSession", "enableSetCookie,omitempty": true, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 157680000, "providerConfig": ""}`
//	conf := new(session.ManagerConfig)
//	if err := json.Unmarshal([]byte(config), conf); err != nil {
//		//t.Fatal("json decode error", err)
//	}
//	globalSessions, _ = session.NewManager("memory", conf)
//	// go globalSessions.GC()
//}
//
//// 处理web页面请求
//func web(rw http.ResponseWriter, req *http.Request) {
//	sess, _ := globalSessions.SessionStart(rw, req)
//	defer sess.SessionRelease(rw)
//	index, _ := viewsIndexHtmlBytes()
//	t, err := template.New("index").Parse(string(index)) //解析模板文件
//	// t, err := template.ParseFiles("web/views/index.html") //解析模板文件
//	if err != nil {
//		logs.Log.Error("%v", err)
//	}
//	//获取pholcus信息
//	data := map[string]interface{}{
//		"title":   config.NAME,
//		"logo":    config.ICON_PNG,
//		"version": config.VERSION,
//		"author":  config.AUTHOR,
//		"mode": map[string]int{
//			"offline": status.OFFLINE,
//			"server":  status.SERVER,
//			"client":  status.CLIENT,
//			"unset":   status.UNSET,
//			"curr":    assetnode.AssetNodeEntity.GetConfig("mode").(int),
//		},
//		"status": map[string]int{
//			"stopped": status.STOPPED,
//			"stop":    status.STOP,
//			"run":     status.RUN,
//			"pause":   status.PAUSE,
//		},
//		"port": assetnode.AssetNodeEntity.GetConfig("port").(int),
//		"ip":   assetnode.AssetNodeEntity.GetConfig("master").(string),
//	}
//	t.Execute(rw, data) //执行模板的merger操作
//}
