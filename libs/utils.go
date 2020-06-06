/**
  create by yy on 2020/5/25
*/

package libs

import (
	"errors"
	"fmt"
	"os"
	"runtime"
)

/**
用于报告错误行数和文件名在哪里，便于找bug
This func is used to report the error line and file name
so that we can find bug quickly.

一般在项目中应用的时候，应该配置一个全局的控制变量，并且打开注释代码块里的注释，
根据你的全局变量进行修改，以达到可以关闭的效果，否则是默认都会报告的
*/
func NewReportError(err error) error {
	// if !config.Config.App.DEBUG {
	//	return err
	// }
	_, fileName, line, _ := runtime.Caller(1)
	data := fmt.Sprintf("%v, report in: %v: in line %v", err, fileName, line)
	return errors.New(data)
}

// 写内容到文件
// params: you can give any type unlimit number data
// Write content to a file.
// If there is not exists, the method will auto create it.
func WriteDataToFile(path string, params ...interface{}) (err error) {
	var (
		ok          bool
		file        *os.File
		content     string
		contentByte []byte
	)

	if file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		// close file stream
		if err = file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if len(params) < 1 {
		fmt.Println(errors.New("null data"))
		return
	}

	// 转化 要写入的内容
	for _, data := range params {
		if content, ok = data.(string); !ok {
			if contentByte, ok = data.([]byte); !ok {
				fmt.Println(errors.New("interface convert to string error"))
				return
			} else {
				_, err = file.Write(contentByte)
			}
		} else {
			_, err = file.Write([]byte(content))
		}

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	return
}
