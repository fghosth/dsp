package util

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

//post传body
func DoBytesPost(url string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*4) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 60)) //设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 60,
			DisableKeepAlives:     true, //关闭保持链接
		},
	}
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		// log.Printf("http.NewRequest,[err=%s][url=%s]", err.Error(), url)
		return []byte(""), err
	}
	//解决"Connection reset by peer"或"EOF"问题
	request.Close = true
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cache-Control", "max-age=0")
	// request.Header.Set("Authorization", "Basic Zmdob3N0aDp6YXExeHN3MkNERSM=")
	var resp *http.Response
	resp, err = client.Do(request)
	// resp, err = http.DefaultClient.Do(request)
	if err != nil {
		// log.Printf("http.Do failed,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// log.Printf("http.Do failed,[err=%s][url=%s]", err, url)
	}
	return b, err
}

//post传body
func DoHTTPGet(url string) ([]byte, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 10)) //设置发送接受数据超时
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 3,
			DisableKeepAlives:     true, //关闭保持链接
		},
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// log.Printf("http.NewRequest,[err=%s][url=%s]", err.Error(), url)
		return []byte(""), err
	}
	//解决"Connection reset by peer"或"EOF"问题
	request.Close = true
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Cache-Control", "max-age=0")
	// request.Header.Set("Authorization", "Basic Zmdob3N0aDp6YXExeHN3MkNERSM=")
	var resp *http.Response
	resp, err = client.Do(request)
	if err != nil {
		// log.Printf("http.Do failed,[err=%s][url=%s]", err, url)
		return []byte(""), err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// log.Printf("http.Do failed,[err=%s][url=%s]", err, url)
	}
	return b, err
}
