package log

import "os"

// 写入文件的路径
type fileLog string

// fileLog 实现了io.Writer中的Write方法
func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}
