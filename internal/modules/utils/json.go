package utils

import (
	"encoding/json"

	"github.com/ouqiang/gocron/internal/modules/logger"
)

// json 格式输出

type response struct {
	Code    int         `json:"code"`    // 状态码 0:成功 非0:失败
	Message string      `json:"message"` // 信息
	Data    interface{} `json:"data"`    // 数据
}

type JsonResponse struct{}

const ResponseSuccess = 0
const ResponseFailure = 1
const UnauthorizedError = 403
const AuthError = 401
const NotFound = 404
const ServerError = 500
const AppNotInstall = 801

const SuccessContent = "操作成功"
const FailureContent = "操作失败"

func JsonResponseByErr(err error) string {
	jsonResp := JsonResponse{}
	if err != nil {
		return jsonResp.CommonFailure(FailureContent, err)
	}

	return jsonResp.Success(SuccessContent, nil)
}

func (j *JsonResponse) Success(message string, data interface{}) string {
	return j.response(ResponseSuccess, message, data)
}

func (j *JsonResponse) Failure(code int, message string) string {
	return j.response(code, message, nil)
}

func (j *JsonResponse) CommonFailure(message string, err ...error) string {
	if len(err) > 0 {
		logger.Warn(err)
	}
	return j.Failure(ResponseFailure, message)
}

func (j *JsonResponse) response(code int, message string, data interface{}) string {
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
