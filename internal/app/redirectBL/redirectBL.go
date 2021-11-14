package redirectBL

import (
	"context"
	"fmt"
	"github.com/a-ivlev/URL-shortener/internal/app/repository/followingBL"
	"github.com/a-ivlev/URL-shortener/internal/app/repository/shortenerBL"
	"time"
)

type Redirect struct {
	shortBL  *shortenerBL.ShortenerBL
	followBL *followingBL.FollowingBL
}

func NewRedirect(shortBL *shortenerBL.ShortenerBL, followBL *followingBL.FollowingBL) *Redirect {
	return &Redirect{
		shortBL:  shortBL,
		followBL: followBL,
	}
}

type Statistic struct {
	ShortLink  string
	TotalCount int
	CreatedAt  time.Time
	FollowList []followingBL.Following
}

func (r *Redirect) CreateShortLink(ctx context.Context, short shortenerBL.Shortener) (*shortenerBL.Shortener, error) {
	createShort, err := r.shortBL.CreateShort(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("create shortener error: %w", err)
	}

	return createShort, nil
}

func (r *Redirect) GetFullLink(ctx context.Context, short shortenerBL.Shortener) (*shortenerBL.Shortener, error) {
	ipaddr := ctx.Value("IP_address").(string)
	getShortener, err := r.shortBL.GetFullLink(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("redirectBL GetFullLink error: %w", err)
	}

	getFollowing, err := r.followBL.SearchFollowing(ctx, getShortener.StatLink, ipaddr)
	if err != nil {
		getFollowing, err = r.followBL.CreateFollowing(ctx, getShortener)
		if err != nil {
			return nil, fmt.Errorf("redirectBL GetFullLink error created new ip-address following list: %w", err)
		}
	}

	getFollowing.IPaddress = ipaddr
	getFollowing.Count += 1
	getFollowing.FollowLinkAt = time.Now()

	getShortener.TotalCount += 1
	updateShort, err := r.shortBL.Update(ctx, *getShortener)
	if err != nil {
		return nil, fmt.Errorf("shortnerBL update error: %w", err)
	}

	_, err = r.followBL.Update(ctx, *getFollowing)
	if err != nil {
		return nil, fmt.Errorf("search following error: %w", err)
	}

	return updateShort, nil
}

func (r *Redirect) GetStatisticList(ctx context.Context, statisticLink string) (*Statistic, error) {
	shortener, err := r.shortBL.ReadShort(ctx, statisticLink)
	if err != nil {
		return nil, fmt.Errorf("get statistic list error: %w", err)
	}

	statistis := &Statistic{
		ShortLink: shortener.ShortLink,
		CreatedAt: shortener.CreatedAt,
		TotalCount: shortener.TotalCount,
	}

	sliceFollowList, err := r.followBL.GetFollowingList(ctx, statisticLink)
	if err != nil {
		return statistis, nil
	}

	statistis.FollowList = sliceFollowList

	return statistis, nil
}
