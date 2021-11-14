package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type testService struct {
	UnimplementedTestServiceServer
}

func main() {
	if tcpServer, err := net.Listen("tcp", "0.0.0.0:8080"); err == nil {
		log.WithFields(log.Fields{
			"addr": tcpServer.Addr(),
		}).Info("tcpServer")
		rpcServer := grpc.NewServer()
		RegisterTestServiceServer(rpcServer, &testService{})
		if err := rpcServer.Serve(tcpServer); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(err)
	}
}

func (testservice *testService) Test(ctx context.Context, req *TestRequest) (*TestReply, error) {
	url := fmt.Sprintf("https://www.rusprofile.ru/search?query=%s", req.Inn)
	log.WithFields(log.Fields{
		"name": req.Inn,
		"url":  url,
	}).Info("Request")
	if resp, err := http.Get(url); err == nil {
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			buf := string(body)
			inn := findField(buf, "id=\"clip_inn\">(.*)<")
			title := findField(buf, "legalName\">(.*)<")
			kpp := findField(buf, "id=\"clip_kpp\">(.*)<")
			director := findDirector(buf)
			return &TestReply{
				Inn:      inn,
				Title:    title,
				Kpp:      kpp,
				Director: director,
			}, nil
		} else {
			log.Error(err)
			return nil, err
		}
	} else {
		log.Error(err)
		return nil, err
	}
}

func findField(body string, pattern string) string {
	regex := regexp.MustCompile(pattern)
	groups := regex.FindAllStringSubmatch(body, -1)
	if len(groups) > 0 {
		return groups[0][1]
	}
	return ""
}

func findDirector(body string) string {
	i := strings.Index(body, "<span class=\"chief-title\">Генеральный директор</span>")
	temp := []rune(body)
	j := strings.Index(string(temp[i:]), "<div>")
	return string(j)
}
