/*
无 format 的图片处理, 只调整大小,不进行格式转换, --- 采用 bild调整大小
*/
package routers

import (
	"fmt"
	"mygo/pics"
	"mygo/utils"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/exp/slog"

	"github.com/gin-gonic/gin"
)

// 创建一个响应对象,用来返回数据格式
type ResponseImageInfo struct {
	File string `json:"file"`
	Msg  string `json:"msg"`
}

var NOFILE = &ResponseImageInfo{
	File: "",
	Msg:  "No File",
}

// 查询参数
type QueryParams struct {
	Sub    string `json:"sub"`    //图像目录下的子目录
	Width  int    `json:"width"`  // 图像宽度, 默认为 0,则是原始大小
	Height int    `json:"height"` // 图像高度, 默认为 0,则是原始大小
	Index  int    `json:"index"`  // 图像索引, 默认为 -1,则是随机
}

// FileUpload godoc
//
//	@Summary		上传文件并保存到服务器
//	@Description	该接口接收一个文件，将其保存到服务器，并将文件信息（包括 MD5 值、文件名、宽高等）存储在内存中
//	@Tags			v1
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file				true	"要上传的文件"
//	@Success		200		{object}	ResponseImageInfo	"Success"
//	@Router			/v1/upload [post]
func FileUpload(ip *utils.ImageProcessor) gin.HandlerFunc {
	imgDir := path.Join(ip.Dir, "upload") // 文件上传的目录(写死)
	return func(c *gin.Context) {
		file, err := c.FormFile("file") // 要和 form-data 中的 name 属性一致
		if err != nil {
			c.JSON(http.StatusBadRequest, NOFILE)
			return
		}
		tempfile, _ := os.CreateTemp("", "*"+file.Filename) // 创建临时文件(采用系统默认的临时目录)
		fmt.Println(" tempfile file:", tempfile.Name())
		defer os.Remove(tempfile.Name()) // 清理临时文件
		defer tempfile.Close()
		// 写入数据到临时文件
		if err := c.SaveUploadedFile(file, tempfile.Name()); err != nil {
			c.JSON(http.StatusInternalServerError, NOFILE)
			return
		}
		dstfinal, _ := RenameWithMd5(tempfile.Name(), imgDir, file.Filename, ip.CalculateMd5)

		if err := c.SaveUploadedFile(file, dstfinal); err != nil {
			c.JSON(http.StatusInternalServerError, NOFILE)
			return
		}
		newInfo, _ := ip.GetSingleImageInfo(dstfinal) //读取新的文件信息
		fmt.Println("dstfinal", dstfinal, "newInfo:", newInfo)
		//添加到内存中
		*ip.SaveDb = append(*ip.SaveDb, newInfo)
		c.JSON(http.StatusOK, ResponseImageInfo{
			File: dstfinal,
			Msg:  "Upload Image Success",
		})

	}
}

// 移动的时候进行重命名
func RenameWithMd5(tempFilePath, imgDir, fileName string, calculateMd5 func(string) (string, error)) (string, error) {
	// 这里可以对文件名进行处理, 返回自己需要的文件名
	// 最后的名字以,文件的 md5 值作为文件名
	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		// 目录不存在，创建目录,  如果path已经是一个目录，MkdirAll不会进行任何操作并返回nil。
		os.MkdirAll(imgDir, os.ModePerm)
	}
	md5, err := calculateMd5(tempFilePath)
	var dstFinal string
	if err != nil {
		slog.Error("计算md5失败", slog.String("err", err.Error()), "file", fileName)
		dstFinal = filepath.Join(imgDir, fileName)
	} else {
		dstFinal = filepath.Join(imgDir, md5+filepath.Ext(fileName)) // 目标路径
	}

	// 移动文件(有坑--原因是跨系统会出问题,特别是 docker和宿主机之间,建议直接写入到目标文件即可)
	// os.Rename(tempFilePath, dstFinal)

	// slog.Info("上传文件:", slog.String("from", tempFilePath), slog.String("to", dstFinal))
	return dstFinal, nil
}

// RandomImage godoc
//
//	@Summary		随机返回一张图片
//	@Description	随机返回一张图片
//	@Tags			v1
//	@Accept			json
//	@Produce		json
//	@Param			sub		query		string				false	"图像目录下的子目录"
//	@Param			width	query		int					false	"图像宽度, 默认为 0,则是原始大小"
//	@Param			height	query		int					false	"图像高度, 默认为 0,则是原始大小"
//	@Param			index	query		int					false	"图像索引, 默认为 -1,则是随机"
//	@Success		200		{object}	ResponseImageInfo	"Success"
//	@Router			/v1/random [get]
func RandomImage(ip *utils.ImageProcessor) gin.HandlerFunc {

	return func(c *gin.Context) {
		// 1. 解析查询参数
		params := ParseQueryParams(c)
		//  params 是 *QueryParams 类型,  在 go 中, 可以直接使用点操作符访问字段:  index := params.Index
		// 								          也可以使用 &params 的方式:   index := (&params).Index
		// 如果params.Sub 为空,则返回源数据
		// 2. 获取图像信息
		var info *utils.ImageInfo
		var path string
		if params.Index >= 0 {
			info = ip.FilterBySub(params.Sub).GetByIndex(params.Index).FirstInfo()
		} else {
			// info = ip.FilterBySub(params.Sub).RandomImageInfo()  //下面等价
			info = ip.RandomImageInfo()
		}
		// 3. 如果查询不到数据，返回404错误
		if info == nil {
			slog.Error("根据条件查询不到数据")
			c.JSON(http.StatusNotFound, NOFILE)
			return
		}
		// 4. 生成文件路径和扩展名
		path = filepath.Join(info.Sub, info.File)
		ext := strings.ToLower(filepath.Ext(path))             //获取扩展名, 带点
		contextType := "image/" + strings.TrimPrefix(ext, ".") //获取扩展名, 并去掉扩展名前的点

		// 5. 检查客户端是否要求返回 JSON 格式
		accept := c.GetHeader("Accept")
		if accept == "application/json" {
			c.JSON(http.StatusOK, ResponseImageInfo{
				File: path,
				Msg:  "success",
			})
			return
		}
		// 6. 读取并处理图像数据(返回二进制数据)
		bytedata, err := pics.LoadImageData(path, ext, params.Width, params.Height)
		if err != nil {
			slog.Error("处理图片时出错", slog.String("err", err.Error()), slog.String("path", path))
			c.JSON(http.StatusInternalServerError, NOFILE)
			return
		}
		// 设置Content-Disposition响应头，指定为附件下载，并为下载文件命名
		// inline：表示文件会在浏览器中打开（如果浏览器支持），但不会自动下载。
		// attachment: 表示以附件形式下载数据。
		// filename="image_title.ext"：指定了文件的名称，浏览器通常会将其作为下载时的默认文件名。
		// 从 path 中提取文件名,build/assets/upload/image.png -> image.png, 采用: filepath.Base 函数即可
		c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="image_%s"`, filepath.Base(path)))
		c.Header("Content-Type", contextType)
		c.Data(http.StatusOK, contextType, bytedata)
	}
}

// RandomImageSrc godoc
//
//	@Summary		随机返回一张图片的地址(以字符串的方式返回,不接受任何参数)
//	@Tags			v1
//	@Accept			plain
//	@Produce		plain
//	@Failure		404		 {string}	 string	 ""
//	@Success		200		 {string}	 "images path"	"Success"
//	@Router			/v1/randomsrc [get]
func RandomImageSrc(ip *utils.ImageProcessor) gin.HandlerFunc {
	// 不支持任何查询操作
	return func(c *gin.Context) {
		info := ip.RandomImageInfo()
		// 3. 如果查询不到数据，返回空
		if info == nil {
			c.String(http.StatusNotFound, "")
			return
		}
		path := filepath.Join(info.Sub, info.File)
		fullpath := filepath.Join(c.Request.Host, path)
		c.String(http.StatusOK, fullpath)
	}
}

// ParseQueryParams 解析查询参数
func ParseQueryParams(c *gin.Context) *QueryParams {
	//获取中间件的值
	// 从上下文中获取设备类型
	deviceType, _ := c.Get("device_type")
	defaultWidth := 0 //

	if deviceType != "pc" {
		defaultWidth = 300 // 移动端默认宽度为300
	}

	params := &QueryParams{
		Sub:    c.DefaultQuery("sub", ""),
		Width:  parseQueryParam(c, "width", defaultWidth),
		Height: parseQueryParam(c, "height", 0),
		Index:  parseQueryParam(c, "index", -1),
	}

	return params
}

// parseQueryParam 是一个辅助函数，用于解析整数查询参数
func parseQueryParam(c *gin.Context, key string, defaultValue int) int {
	value, err := strconv.Atoi(c.DefaultQuery(key, strconv.Itoa(defaultValue)))
	if err != nil {
		return defaultValue
	}
	return value
}
