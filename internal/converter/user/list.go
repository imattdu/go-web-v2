package user

import (
	"context"
	"reflect"

	"github.com/imattdu/go-web-v2/internal/dto/user"
	userAPI "github.com/imattdu/go-web-v2/internal/service/user/api"

	"github.com/jinzhu/copier"
)

func ListRequestToParams(ctx context.Context, request user.ListRequest) (params userAPI.ListParams, err error) {
	err = copier.CopyWithOption(&params, &request, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
	return params, err
}

func ListParamsToResponse(ctx context.Context, result userAPI.ListResult) (response user.ListResponse, err error) {
	err = copier.CopyWithOption(&response, &result, copier.Option{
		IgnoreEmpty:   true, // 空值不覆盖
		DeepCopy:      true, // 深拷贝
		CaseSensitive: true, // 不忽略大小写
		//Converters: []copier.TypeConverter{
		//	{
		//		SrcType: reflect.TypeOf(""), // 来源类型
		//		DstType: reflect.TypeOf(0),  // 目标类型
		//		Fn: func(src interface{}) (interface{}, error) {
		//			return len(src.(string)), nil // 转换逻辑
		//		},
		//	},
		//},
		FieldNameMapping: []copier.FieldNameMapping{
			{
				SrcType: reflect.TypeOf(""),
				DstType: reflect.TypeOf(""),
				Mapping: map[string]string{
					"username": "Username",
				},
			},
		},
	})
	return response, err
}
