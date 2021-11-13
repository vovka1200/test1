package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type testService struct {
	UnimplementedTestServiceServer
}

func (testservice *testService) Test(ctx context.Context, req *TestRequest) (*TestReply, error) {
	log.WithFields(log.Fields{
		"name": req.Inn,
	}).Info("Request")
	return &TestReply{Inn: "Test"}, nil
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
