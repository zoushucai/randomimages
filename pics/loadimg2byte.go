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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/transform"
	"golang.org/x/image/webp"
)

// LoadImageData 处理图像文件，返回其二进制数据
func LoadImageData(path, ext string, width, height int) ([]byte, error) {
	if ext == ".webp" {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		return io.ReadAll(file)
	}
	// 处理非 webp 格式的图片
	// 读取图像
	img, err := ReadImage(path)
	if err != nil {
		return nil, fmt.Errorf("读取图片失败: %w", err)
	}

	// 调整图像大小
	newImg, err := ImageResize(img, width, height)
	if err != nil {
		return nil, fmt.Errorf("图片缩放失败: %w", err)
	}

	// 转换图像为字节数据
	return ConvertImageToBytes(newImg, ext)
}

func ReadImage(path string) (image.Image, error) {
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
	case ".webp":
		img, err = webp.Decode(imgFile)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return nil, err
	}
	return img, nil
}

// 调整图像大小
// img image.Image 输入的图像, 以及宽高
func ImageResize(img image.Image, width int, height int) (newImage image.Image, err error) {

	if width == 0 && height == 0 {
		width = img.Bounds().Dx()
		height = img.Bounds().Dy()
	} else if width == 0 {
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
	// newImage = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	return newImage, nil
	// convertedImage, _ := ConvertImageToBytes(newImage, ext)
	// return convertedImage, nil
}

// 将图像转换为字节数组
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
