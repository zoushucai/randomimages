package pics

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/webp"
)

// converts a WebP image to a specified format (e.g., "png" or "jpeg").
func ConvertWebPToFormat(input, format string) (outfile string, err error) {
	// 打开 WebP 文件
	f, err := os.Open(input)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// 解码 WebP 图像
	img, err := webp.Decode(f)
	if err != nil {
		return "", fmt.Errorf("解码 WebP 图像失败: %w", err)
	}

	// 根据指定格式保存图像

	switch format {
	case "png":
		// 把原始的 webp 扩展名给替换为 png
		outfile := strings.Replace(input, ".webp", ".png", -1)
		return outfile, saveAsPNG(outfile, img)
	case "jpeg":
		outfile := strings.Replace(input, ".webp", ".jpeg", -1)
		return outfile, saveAsJPEG(outfile, img)
	default:
		return outfile, fmt.Errorf("不支持的格式: %s", format)
	}
}

// saveAsPNG saves the given image as a PNG file.
func saveAsPNG(filename string, img image.Image) error {
	pngFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("无法创建 PNG 文件: %w", err)
	}
	defer pngFile.Close()

	if err := png.Encode(pngFile, img); err != nil {
		return fmt.Errorf("PNG 编码失败: %w", err)
	}

	return nil
}

// saveAsJPEG saves the given image as a JPEG file with a specified quality.
func saveAsJPEG(filename string, img image.Image) error {
	jpegFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("无法创建 JPEG 文件: %w", err)
	}
	defer jpegFile.Close()

	if err := jpeg.Encode(jpegFile, img, &jpeg.Options{Quality: 90}); err != nil {
		return fmt.Errorf("JPEG 编码失败: %w", err)
	}

	return nil
}
