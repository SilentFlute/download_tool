package server

import (
	"net/http"
	"strconv"
	"log"
	"download_tool/conf"
	"net"
	"regexp"
	"strings"
	"errors"
)

var fileServer http.Handler
var contentMap = map[string]string {
	"css":  "text/css; charset=utf-8",
	"gif":  "image/gif",
	"htm":  "text/html; charset=utf-8",
	"html": "text/html; charset=utf-8",
	"jpg":  "image/jpeg",
	"mp3":  "audio/mp3",
	"mp4":  "video/mpeg4",
	"js":   "application/javascript",
	"wasm": "application/wasm",
	"pdf":  "application/pdf",
	"png":  "image/png",
	"svg":  "image/svg+xml",
	"xml":  "text/xml; charset=utf-8",
	"m3u8":  "video/mpeg4",
}

func handleContentType(w http.ResponseWriter, req *http.Request) {
	reg, err := regexp.Compile("^.*\\.(.+)$")
	if err != nil {
		log.Fatalln("正则编译错")
		return
	}

	uri := req.RequestURI
	if !reg.MatchString(uri) {
		return
	}

	strs := reg.FindAllStringSubmatch(uri, -1)
	suffix := strings.TrimSpace(strs[0][1])
	if ctp, ok := contentMap[suffix]; ok {
		contentType := make([]string, 1)
		contentType[0] = ctp
		w.Header()["Content-Type"] = contentType
	}
}

func handleTerminate(w http.ResponseWriter, req *http.Request) {
	log.Println(req.RequestURI)
	if req.RequestURI == "/download" {
		err := req.ParseForm()
		if err != nil {
			w.WriteHeader(500)
			return
		}

		url := req.Form.Get("url")
		name := req.Form.Get("name")
		cookie := req.Form.Get("cookie")
		referer := req.Form.Get("referer")
		err = download(name, referer, cookie, url)
		w.WriteHeader(200)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte("下载成功"))
		}
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	log.Println("requesuri is", req.RequestURI)
	handleContentType(w, req)

	fileServer.ServeHTTP(w, req);
}

func Start()  {
	ps := conf.GetConfigValue("port")
	var port = 8080;
	var configPort = -1

	if ps != "" {
		log.Println("将使用端口[", ps, "]启动服务")
		p, err := strconv.Atoi(ps)
		if err != nil {
			log.Fatalln("端口配置错误:" + ps)
			return
		}

		configPort = p
	}

	if configPort != -1 {
		port = configPort
	}

	current := conf.GetCurrentDirectory();
	dir := current + "/movie";
	fileServer = http.FileServer(http.Dir(dir))

	log.Println(dir)

	http.Handle("/download", http.HandlerFunc(handleTerminate))
	http.Handle("/", http.HandlerFunc(handler))

	log.Println("handle...")

	err := listen(port)
	log.Println("err:", err == nil, port)

	for err != nil && configPort == -1 {
		port++
		err = listen(port)
		log.Println("err:", err == nil, port)
		if(port > 60000) {
			break
		}
	}

	if err != nil {
		log.Fatalln("端口[", port, "]已经被占用，无法启动，请在config.ini中修改端口号，并重启")
	}
}

func listen(port int) error {
	conn, err := net.Dial("tcp", "localhost:" + strconv.Itoa(port))
	if err == nil {
		conn.Close()
		return errors.New("端口被占用");
	}

	return http.ListenAndServe(":" + strconv.Itoa(port), nil)
}
