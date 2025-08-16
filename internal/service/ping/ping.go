package ping

import (
	"context"
	"fmt"
	"github.com/imattdu/go-web-v2/internal/database/redis"
	"time"
)

type Service interface {
	Ping(ctx context.Context) string
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Ping(ctx context.Context) string {
	//_ = httpclientresty.Post(ctx, &httpclientresty.Req{
	//	Service: errorx.ServiceMysql,
	//	Meta: httpclientresty.ReqMeta{
	//		URL: "http://127.0.0.1:8001/go-web-v2/user/list",
	//		RequestBody: user.ListRequest{
	//			Username: "matt",
	//		},
	//		Timeout: time.Millisecond * 10,
	//		//ResponseBody: &user.ListResponse{},
	//		//IsError: func(response *http.Response) error {
	//		//	return errorx.New(errorx.ErrOptions{
	//		//		ErrMeta: errorx.ErrMeta{
	//		//			ErrType: errorx.ErrTypeBiz,
	//		//		},
	//		//		Err: errors.New("abc"),
	//		//	})
	//		//},
	//	},
	//})

	var ss string
	ss = "hh"
	rsp, err := redis.TTL(ctx, &redis.KVEntry{
		KeyPre:  "f99",
		VBody:   &ss,
		BaseTTL: time.Minute * 3,
	})
	fmt.Println(rsp, ss)
	if err != nil {
		fmt.Println(err.Error())
	}
	return "pong"
}
