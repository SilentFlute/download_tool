package server

import (
	"net/http"
	"time"
	"os"
	"download_tool/conf"
	"io/ioutil"
	"strings"
	"container/list"
	"log"
	"errors"
)

type task struct {
	fileName string
	ts string
	err chan error
}

func download(name string, referer string, cookie string, url string) error {
	index, err := downloadIndex(name, referer, cookie, url);
	if err != nil {
		return err;
	}

	taskList := list.New()
	arr := strings.Split(index, "\n");
	for _, val := range arr {
		if len(val)== 0 || val[0:1] == "#" {
			continue
		}

		fileName := conf.GetCurrentDirectory() + "/movie/" + name + "/" + val
		ts := url[0:strings.LastIndex(url, "/")] + "/" + val;

		t := &task{
			fileName:fileName,
			ts:ts,
			err:make(chan error),
		}

		taskList.PushBack(t)
	}

	defer func() {
		for e:=taskList.Front(); e != nil; e = e.Next() {
			t := e.Value.(*task)
			close(t.err)
		}
	}()

	for e:=taskList.Front(); e != nil; e = e.Next() {
		t := e.Value.(*task)

		go func() {
			t.err<-downloadTs(t.fileName, referer, cookie, t.ts)
		}()
	}

	return retryDownload(0, taskList, referer, cookie)
}

func retryDownload(idx int, taskList *list.List, referer string, cookie string) error {
	flag := true
	errList := list.New()
	for e:=taskList.Front(); e != nil; e = e.Next() {
		t := e.Value.(*task)
		if <-t.err != nil {
			flag = false
			errList.PushBack(e.Value)
			go func() {
				t.err <- downloadTs(t.fileName, referer, cookie, t.ts)
			}()
		}
	}

	if flag {
		log.Println("下载成功")
		return nil
	}

	if idx >= 10 {
		log.Fatalf("有部分分片文件下载失败")
		return errors.New("部分文件下载失败")
	}

	return retryDownload(idx + 1, errList, referer, cookie)
}

func downloadTs(fileName string, referer string, cookie string, url string) error {
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	client := &http.Client{}
	//50s超时
	var s time.Duration = 50000000000
	client.Timeout = s

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err;
	}

	req.Header.Add("User-Agent", userAgent)

	if referer != "" {
		req.Header.Add("Referer", referer)
	}

	if cookie != "" {
		req.Header.Add("Cookie", cookie)
	}

	response, err := client.Do(req)
	if err != nil {
		return err;
	}

	defer response.Body.Close()

	f, err := os.Create(fileName)
	if err != nil {
		return err;
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err;
	}

	_, err = f.Write(body)

	return err
}

func downloadIndex(name string, referer string, cookie string, url string) (string, error) {
	userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	client := &http.Client{}
	//50s超时
	var s time.Duration = 50000000000
	client.Timeout = s

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err;
	}

	req.Header.Add("User-Agent", userAgent)

	if referer != "" {
		req.Header.Add("Referer", referer)
	}

	if cookie != "" {
		req.Header.Add("Cookie", cookie)
	}

	response, err := client.Do(req)
	if err != nil {
		return "", err;
	}

	defer response.Body.Close()

	dir := conf.GetCurrentDirectory() + "/movie/" + name
	err = os.MkdirAll(dir, 0770)
	if err != nil {
		return "", err;
	}

	f, err := os.Create(dir + "/index.m3u8")
	if err != nil {
		return "", err;
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err;
	}

	_, err = f.Write(body)

	return string(body), err
}
