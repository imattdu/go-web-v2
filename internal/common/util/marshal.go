package util

import (
	"context"
	"encoding/json"
	logger2 "github.com/imattdu/go-web-v2/internal/common/util/logger"
)

func Marshal(ctx context.Context, body interface{}) (string, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		logger2.Warn(ctx, logger2.TagUndef, map[string]interface{}{
			"msg":  "util.Marshal failed",
			"err":  err.Error(),
			"body": body,
		})
		return "", err
	}
	return string(bodyBytes), nil
}

func Unmarshal(ctx context.Context, str string, body interface{}) error {
	if err := json.Unmarshal([]byte(str), &body); err != nil {
		logger2.Warn(ctx, logger2.TagUndef, map[string]interface{}{
			"msg":  "util.Unmarshal failed",
			"err":  err.Error(),
			"str":  str,
			"body": body,
		})
		return err
	}
	return nil
}
