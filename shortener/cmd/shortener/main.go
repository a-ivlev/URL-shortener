package main

import (
	"CourseProjectBackendDevGoLevel-1/shortener/internal/api/chiRouter"
	"CourseProjectBackendDevGoLevel-1/shortener/internal/api/handler"
	"CourseProjectBackendDevGoLevel-1/shortener/internal/api/server"
	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/redirectBL"
	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/repository/followingBL"
	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/repository/shortenerBL"
	"CourseProjectBackendDevGoLevel-1/shortener/internal/app/starter"
	"CourseProjectBackendDevGoLevel-1/shortener/internal/db/inmemoryDB"
	"context"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	srvPort := os.Getenv("SHORT_SRV_PORT")
	if srvPort == "" {
		srvPort = "8035"
	}

	//cliHost := os.Getenv("SHORT_CLI_HOST")
	//if cliHost == "" {
	//	cliHost = "localhost"
	//}
	//ctx = context.WithValue(ctx, "cliHost", cliHost)

	shortdb := inmemoryDB.NewShortenerMapDB()
	followdb := inmemoryDB.NewFollowingMapDB()

	shortBL := shortenerBL.NewShotenerBL(shortdb)
	followBL := followingBL.NewFollowingBL(followdb)
	redirBL := redirectBL.NewRedirect(shortBL, followBL)

	app := starter.NewApp(redirBL)
	handlers := handler.NewHandlers(redirBL)
	chi := chiRouter.NewChiRouter(handlers)
	srv := server.NewServer(":"+srvPort, chi)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go app.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
