package fsExtra

import (
	"io/fs"
	"os"
)

// WriteStringFile 写入文件
// param:filePath 文件路径
// param:content 内容
func WriteStringFile(filePath string, content string) error {
	return os.WriteFile(filePath, []byte(content), fs.ModePerm)
}
