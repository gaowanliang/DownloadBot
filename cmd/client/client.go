package client

import (
	pb "DownloadBot/api/DownloadBot/v1"
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/internal/config"
	logger "DownloadBot/tool/zap"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type client struct {
	pb pb.GreeterClient
}

var Method = new(client)

func (c client) InitClient() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", config.GetServerIP(), config.GetServerPort()), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(i18nLoc.LocText("did not connect: %v"), err)
	}
	// defer conn.Close()
	clientPB := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//c.ctx, _ = context.WithTimeout(context.Background(), time.Second)
	r, err := clientPB.SayHello(ctx, &pb.HelloRequest{Name: config.GetSign()})
	if err != nil {
		logger.Panic(i18nLoc.LocText("connect server failed: %v"), err)
	}
	logger.Info(i18nLoc.LocText("connected server: %v"), r.GetMessage())
	Method.pb = clientPB
	go func() {
		for {
			// ensure client is alive
			logger.Debug("Ping Send: %v", config.GetSign())
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_, err := Method.pb.Ping(ctx, &pb.Ping{ClientName: config.GetSign()})
			if err != nil {
				logger.Panic("Ping failed: %v", err)
			}
			time.Sleep(time.Second * 20)
		}
	}()

}

func (c client) SendSuddenMessage(msg string) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := c.pb.SendSuddenMessage(ctx, &pb.SuddenMessage{Message: msg})
	if err != nil {
		logger.Error(i18nLoc.LocText("could not send sudden message: %v"), err)
	}
}

func (c client) TMStop(gid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := c.pb.TMStop(ctx, &pb.TMStopMsg{Gid: gid})
	if err != nil {
		logger.Error(i18nLoc.LocText("could not send TMStop: %v"), err)
	}
}
