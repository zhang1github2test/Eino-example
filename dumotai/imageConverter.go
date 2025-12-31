package main

import "Eino-example/util"

// 定义图像转换接口
type ImageConverter interface {
	ImageToBase64(path string) (mimeType string, data string, err error)
}

// 默认实现
type DefaultImageConverter struct{}

func (c *DefaultImageConverter) ImageToBase64(path string) (mimeType string, data string, err error) {
	return util.ImageToBase64(path)
}

// 全局变量用于测试时替换
var imageConverter ImageConverter = &DefaultImageConverter{}
