package ping

import (
	"context"
)

type Service interface {
	Ping(ctx context.Context) string
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Ping(ctx context.Context) string {
	//_ = httpclient.Post(ctx, &httpclient.Req{
	//	Service: errorx.ServiceMysql,
	//	Meta: httpclient.ReqMeta{
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
	return "pong"
}
