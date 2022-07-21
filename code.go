package echotool

import (
	"net/http"
)

const (
	CodeOKZero         = 0
	CodeOK             = 20000
	CodeCreated        = 20100
	CodePartialContent = 20600

	CodeMultipleChoices   = 30000
	CodeMovedPermanently  = 30100
	CodeFound             = 30200
	CodeSeeOther          = 30300
	CodeNotModified       = 30400
	CodeUseProxy          = 30500
	CodeTemporaryRedirect = 30700
	CodePermanentRedirect = 30800

	CodeBadRequest      = 40000
	CodeUnauthorized    = 40100
	CodeForbidden       = 40300
	CodeNotFound        = 40400
	CodeTooManyRequests = 42900
	CodeValidateErr     = 45000

	CodeInternalErr        = 50000
	CodeServiceUnavailable = 50300
	CodeParseDataErr       = 50013
	CodeMySQLErr           = 50014
	CodePostgreSQLErr      = 50015
	CodeRedisErr           = 50016
	CodeClickHouseErr      = 50017
	CodeNSQErr             = 50030
	CodeKafkaErr           = 50031
	CodeRocketMQErr        = 50032
	CodeBindErr            = 50040
	CodeEncodeErr          = 50041
	CodeDownstreamErr      = 50099
)

var codeMsg = map[int]string{
	CodeOKZero:         "success",
	CodeOK:             "success",
	CodeCreated:        "created",
	CodePartialContent: "partial content",

	CodeMultipleChoices:   "multiple choices",
	CodeMovedPermanently:  "moved permanently",
	CodeFound:             "found",
	CodeSeeOther:          "see other",
	CodeNotModified:       "not modified",
	CodeUseProxy:          "use proxy",
	CodeTemporaryRedirect: "temporary redirect",
	CodePermanentRedirect: "permanent redirect",

	CodeBadRequest:      "bad request",
	CodeUnauthorized:    "unauthorized",
	CodeForbidden:       "forbidden",
	CodeNotFound:        "not found",
	CodeTooManyRequests: "too many requests",
	CodeValidateErr:     "validate error",

	CodeInternalErr:        "internal error",
	CodeServiceUnavailable: "service unavailable",
	CodeParseDataErr:       "parse request data error",
	CodeMySQLErr:           "mysql error",
	CodePostgreSQLErr:      "postgresql error",
	CodeRedisErr:           "redis error",
	CodeClickHouseErr:      "clickhouse error",
	CodeNSQErr:             "nsq error",
	CodeKafkaErr:           "kafka error",
	CodeRocketMQErr:        "rocketmq error",
	CodeBindErr:            "bind error",
	CodeEncodeErr:          "encode error",
	CodeDownstreamErr:      "downstream error",
}

func CodeMsg(code int) string {
	if msg, exists := codeMsg[code]; exists {
		return msg
	}
	return "unknown code"
}

var httpStatus = map[int]int{
	CodeOKZero:         http.StatusOK,
	CodeOK:             http.StatusOK,
	CodeCreated:        http.StatusCreated,
	CodePartialContent: http.StatusPartialContent,

	CodeMultipleChoices:   http.StatusMultipleChoices,
	CodeMovedPermanently:  http.StatusMovedPermanently,
	CodeFound:             http.StatusFound,
	CodeSeeOther:          http.StatusSeeOther,
	CodeNotModified:       http.StatusNotModified,
	CodeUseProxy:          http.StatusUseProxy,
	CodeTemporaryRedirect: http.StatusTemporaryRedirect,
	CodePermanentRedirect: http.StatusPermanentRedirect,

	CodeBadRequest:      http.StatusBadRequest,
	CodeUnauthorized:    http.StatusUnauthorized,
	CodeForbidden:       http.StatusForbidden,
	CodeNotFound:        http.StatusNotFound,
	CodeTooManyRequests: http.StatusTooManyRequests,
	CodeValidateErr:     http.StatusBadRequest,

	CodeInternalErr:        http.StatusInternalServerError,
	CodeServiceUnavailable: http.StatusInternalServerError,
	CodeParseDataErr:       http.StatusInternalServerError,
	CodeMySQLErr:           http.StatusInternalServerError,
	CodePostgreSQLErr:      http.StatusInternalServerError,
	CodeRedisErr:           http.StatusInternalServerError,
	CodeClickHouseErr:      http.StatusInternalServerError,
	CodeNSQErr:             http.StatusInternalServerError,
	CodeKafkaErr:           http.StatusInternalServerError,
	CodeRocketMQErr:        http.StatusInternalServerError,
	CodeBindErr:            http.StatusInternalServerError,
	CodeEncodeErr:          http.StatusInternalServerError,
	CodeDownstreamErr:      http.StatusInternalServerError,
}

const (
	UnknownStatus = 999
)

func HTTPStatus(code int) int {
	if status, exists := httpStatus[code]; exists {
		return status
	}
	return UnknownStatus
}

// RegisterCode will not cover code and status which exists.
func RegisterCode(code int, msg string, status int) bool {
	if _, exists := codeMsg[code]; exists {
		return false
	}
	if _, exists := httpStatus[code]; exists {
		return false
	}

	codeMsg[code] = msg
	httpStatus[code] = status
	return true
}

// ForceRegisterCode will force to cover codeMsg map and httpStatus map.
func ForceRegisterCode(code int, msg string, status int) {
	codeMsg[code] = msg
	httpStatus[code] = status
}
