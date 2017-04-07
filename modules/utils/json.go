package utils

import (
    "encoding/json"
    "github.com/ouqiang/gocron/modules/logger"
)

// json 格式输出

type response struct {
    Code    int         `json:"code"`    // 状态码 0:成功 非0:失败
    Message string      `json:"message"` // 信息
    Data    interface{} `json:"data"`    // 数据
}

type Json struct{}

const ResponseSuccess = 0
const ResponseFailure = 1

func (j *Json) Success(message string, data interface{}) string {
    return j.response(ResponseSuccess, message, data)
}

func (j *Json) Failure(code int, message string) string {
    return j.response(code, message, nil)
}

func (j *Json) response(code int, message string, data interface{}) string {
    resp := response{
        Code:    code,
        Message: message,
        Data:    data,
    }

    result, err := json.Marshal(resp)
    if err != nil {
        logger.Error(err)
    }

    return string(result)
}
