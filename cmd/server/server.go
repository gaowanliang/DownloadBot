package server

import (
	pb "DownloadBot/api/DownloadBot/v1"
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/internal/config"
	"DownloadBot/internal/server/clientManage"
	"DownloadBot/tool/output/telegram"
	logger "DownloadBot/tool/zap"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"time"
)

// server is used to implement DownloadBot.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements DownloadBot.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if !clientManage.ClientSearch(in.GetName()) {
		clientManage.ClientList = append(clientManage.ClientList, clientManage.ClientInfo{ClientName: in.GetName(), AlivePoint: 1})
	}
	logger.Info(i18nLoc.LocText("Client Log in: %v"), in.GetName())
	return &pb.HelloReply{Message: config.GetSign()}, nil
}

func (s *server) SendSuddenMessage(ctx context.Context, in *pb.SuddenMessage) (*pb.Status, error) {
	logger.Debug("SuddenMessage: %v", in.GetMessage())
	telegram.SendTelegramSuddenMessage(in.GetMessage())
	return &pb.Status{Code: 0, Message: ""}, nil
}

func (s *server) TMStop(ctx context.Context, in *pb.TMStopMsg) (*pb.Status, error) {
	logger.Debug("TMStop: %v", in.GetGid())
	telegram.TMSelectMessageChan <- in.GetGid()
	return &pb.Status{Code: 0, Message: ""}, nil
}

func (s *server) Ping(ctx context.Context, in *pb.Ping) (*pb.Status, error) {
	logger.Debug("Ping received from 『%v』", in.GetClientName())
	if !clientManage.ClientSearch(in.GetClientName()) {
		clientManage.ClientList = append(clientManage.ClientList, clientManage.ClientInfo{ClientName: in.GetClientName(), AlivePoint: 1})
		logger.Info(i18nLoc.LocText("Client Log in: %v"), in.GetClientName())
	} else {
		for i, v := range clientManage.ClientList {
			if v.ClientName == in.GetClientName() {
				clientManage.ClientList[i].AlivePoint = 1
			}
		}
	}
	return &pb.Status{Code: 0, Message: ""}, nil
}

func clientAlive() {
	logger.Debug("clientAlive has been launched")
	for {
		for i, v := range clientManage.ClientList {
			if v.AlivePoint == 0 {
				// delete client
				clientManage.ClientList = append(clientManage.ClientList[:i], clientManage.ClientList[i+1:]...)
				logger.Info(i18nLoc.LocText("Client 『%v』 is offline"), v.ClientName)
				// telegram.SuddenMessageChan <- fmt.Sprintf(i18nLoc.LocText("Client 『%v』 is offline"), v.ClientName)
			} else {
				clientManage.ClientList[i].AlivePoint = 0
				// logger.Debug("Client %v maybe dead", v.clientName)
			}
		}
		time.Sleep(time.Second * 30)
	}

}

func StartServer() {
	var err error
	var lis net.Listener

	lis, err = net.Listen("tcp", fmt.Sprintf(":%d", config.GetServerPort()))
	logger.DropErr(err)

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	logger.Info(i18nLoc.LocText("server listening at %v"), lis.Addr())

	go clientAlive()

	if err := s.Serve(lis); err != nil {
		logger.DropErr(err)
	}

}
