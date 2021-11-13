package starter

import (
	"github.com/a-ivlev/URL-shortener/shortener/internal/app/redirectBL"
	"context"
	"sync"
)

type App struct {
	redirectBL *redirectBL.Redirect
}

func NewApp(redirectBL *redirectBL.Redirect) *App {
	app := &App{
		redirectBL: redirectBL,
	}
	return app
}


type APIServer interface {
	Start(redirectBL *redirectBL.Redirect)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start(a.redirectBL)
	<-ctx.Done()
	hs.Stop()
}
