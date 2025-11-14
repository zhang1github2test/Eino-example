package util

import (
	"encoding/base64"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// ImageToBase64 读取图片文件，返回 MIME 类型和不带头的 Base64 编码内容
func ImageToBase64(imagePath string) (mimeType string, base64Data string, err error) {
	// 读取图片文件内容
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", "", err
	}

	// 检测 MIME 类型
	// http.DetectContentType 使用前 512 个字节来判断类型
	// 对于某些文件，如果魔数不够，可能需要手动根据扩展名判断
	mimeType = http.DetectContentType(data)

	// 如果 DetectContentType 无法识别，可以尝试根据文件扩展名来补充
	if mimeType == "application/octet-stream" {
		// 尝试从文件扩展名推断
		ext := filepath.Ext(imagePath)
		if mt := mime.TypeByExtension(ext); mt != "" {
			mimeType = mt
		}
	}

	// 将字节数组编码为 Base64 字符串
	base64Data = base64.StdEncoding.EncodeToString(data)

	return mimeType, base64Data, nil
}
