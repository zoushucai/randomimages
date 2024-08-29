package pics

/*
简单的图像处理
*/
import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/transform"
)

// ProcessImage 读取图像，进行缩放，并返回图像的字节切片
func ProcessImage(path string, width, height int) ([]byte, error) {
	// 打开图像文件
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(path))

	// 解码图像
	var img image.Image
	switch ext {
	case ".jpeg", ".jpg", ".png", ".gif":
		img, _, err = image.Decode(imgFile)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return nil, err
	}

	// 缩放图像
	var newImage image.Image
	if width == 0 && height == 0 {
		newImage = img
	} else {
		if width == 0 {
			// 如果只指定高度，按比例计算宽度
			ratio := float64(height) / float64(img.Bounds().Dy())
			width = int(ratio * float64(img.Bounds().Dx()))
		} else if height == 0 {
			// 如果只指定宽度，按比例计算高度
			ratio := float64(width) / float64(img.Bounds().Dx())
			height = int(ratio * float64(img.Bounds().Dy()))
		}
		// NearestNeighbor
		newImage = transform.Resize(img, width, height, transform.Lanczos)

	}
	convertedImage, _ := ConvertImageToBytes(newImage, ext)
	return convertedImage, nil
}

func ConvertImageToBytes(img image.Image, format string) ([]byte, error) {
	var buffer bytes.Buffer
	var err error

	switch format {
	case ".jpeg":
		err = jpeg.Encode(&buffer, img, nil)
	case ".png":
		err = png.Encode(&buffer, img)
	case ".gif":
		err = gif.Encode(&buffer, img, nil)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
