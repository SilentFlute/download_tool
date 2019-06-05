package conf

import (
	"os"
	"io/ioutil"
	"log"
	"strings"
	"regexp"
	"syscall"
	"path/filepath"
)

var data map[string]string

func GetConfigValue(key string) string {
	return data[key]
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func init() {
	reg, err := regexp.Compile("(.*?)=(.*)")
	if err != nil {
		log.Fatalln("正则编译错")
		return
	}

	reg2, err := regexp.Compile("^.*\\[.+\\].*$")
	if err != nil {
		log.Fatalln("正则编译错")
		return
	}

	f, err := os.Open(GetCurrentDirectory() + "/config.ini")
	if err != nil {
		log.Println("未提供配置文件或打开配置文件出错:", err.Error())
		return
	}

	defer f.Close()

	bytes, err := ioutil.ReadAll(f)

	if err != nil {
		log.Fatalln("读取配置文件出错:" + err.Error())
		return
	}

	configStr := string(bytes)
	configStr = strings.Replace(configStr, "\r", "", -1)
	if strings.TrimSpace(configStr) == "" {
		return
	}

	lines := strings.Split(configStr, "\n")

	data = make(map[string]string)

	for i := range lines {
		line := lines[i]
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.Index(line, "#") == 0 {
			continue
		}

		if reg2.MatchString(line) && !reg.MatchString(line) {
			continue
		}

		if !reg.MatchString(line) {
			log.Fatalln("配置格式不符:" + line + " asc:", []byte(line)[2])
			syscall.Exit(-1)
		}

		strs := reg.FindAllStringSubmatch(line, -1)
		key := strings.TrimSpace(strs[0][1])
		val := strings.TrimSpace(strs[0][2])

		data[key] = val
 	}

 	log.Println("加载配置文件成功")
}
