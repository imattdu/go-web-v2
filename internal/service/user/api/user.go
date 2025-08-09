package api

import (
	"context"
)

type Service interface {
	List(ctx context.Context, params ListParams) (ListResult, error)
}
