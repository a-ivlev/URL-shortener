package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/a-ivlev/URL-shortener/shortener/internal/api/chiRouter"
	"github.com/a-ivlev/URL-shortener/shortener/internal/api/handler"
	"github.com/a-ivlev/URL-shortener/shortener/internal/api/server"
	"github.com/a-ivlev/URL-shortener/shortener/internal/app/redirectBL"
	"github.com/a-ivlev/URL-shortener/shortener/internal/app/repository/followingBL"
	"github.com/a-ivlev/URL-shortener/shortener/internal/app/repository/shortenerBL"
	"github.com/a-ivlev/URL-shortener/shortener/internal/app/starter"
	"github.com/a-ivlev/URL-shortener/shortener/internal/db/inmemoryDB"
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
