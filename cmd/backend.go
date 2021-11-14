package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/a-ivlev/URL-shortener/internal/api/chiRouter"
	"github.com/a-ivlev/URL-shortener/internal/api/handler"
	"github.com/a-ivlev/URL-shortener/internal/api/server"
	"github.com/a-ivlev/URL-shortener/internal/app/redirectBL"
	"github.com/a-ivlev/URL-shortener/internal/app/repository/followingBL"
	"github.com/a-ivlev/URL-shortener/internal/app/repository/shortenerBL"
	"github.com/a-ivlev/URL-shortener/internal/app/starter"
	"github.com/a-ivlev/URL-shortener/internal/db/inmemoryDB"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	srvPort := os.Getenv("PORT")
	if srvPort == "" {
		log.Fatal("unknown PORT = ", srvPort)
	}

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
