package main

import (
	"context"
	"fmt"
	"github.com/imattdu/go-web-v2/internal/common/cctx"
)

type A struct {
	Name string
}

func t1(ctx context.Context) {
	cctx.Set(ctx, "cc", A{
		Name: "t1",
	})
}

func main() {
	ctx := cctx.NewContext(context.Background(), map[string]any{
		"a": 123,
		"b": "haha",
	})
	t1(ctx)

	//data := cctx.GetData(ctx)
	fmt.Println(cctx.Get(ctx, "hahah"))
}
