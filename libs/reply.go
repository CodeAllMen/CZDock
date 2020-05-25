/**
  create by yy on 2019-07-29
*/

package libs

type Reply struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Reply {
	return &Reply{
		Code: 0,
		Msg:  "",
		Data: data,
	}
}

func Error(msg string) *Reply {
	return &Reply{
		Code: 1,
		Msg:  msg,
		Data: "",
	}
}

func CustomReply(code int, msg string, data ...interface{}) *Reply {
	var (
		replyData interface{}
	)

	if len(data) > 0 {
		replyData = data[0]
	} else {
		replyData = ""
	}

	return &Reply{
		Code: code,
		Msg:  msg,
		Data: replyData,
	}

}
