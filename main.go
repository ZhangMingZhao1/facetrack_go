package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
	"gocv.io/x/gocv"
)

//运行方式:go run human.go 0 haarcascade_frontalface_default.xml
func main() {
	// parse args
	//deviceID, _ := strconv.Atoi(os.Args[1])
	//获取int类型值0
	deviceID, _ := strconv.Atoi("0")
	//xmlFile := os.Args[2]
	xmlFile := "haarcascade_frontalface_default.xml"

	// open webcam 开启视频
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window 并设定名称
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix 图片矩阵
	img := gocv.NewMat()
	defer img.Close()

	//加载分类器以识别人脸
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	//分类器加载图像xml
	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	//开启循环，读取
	for {
		//img为矩阵，将图片存储在矩阵中
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		//如果读取的图片不为空
		if img.Empty() {
			continue
		}
		//DetectMultiscale检测输入mat图像中的人脸区域,如果有多张脸,以切片的形式返回脸占用的矩形区域
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))
		//读取timg.png这张狗头图片
		file, _ := os.Open("timg.png")
		//将狗头图片转换成img.Image类型
		srcPng, _ := png.Decode(file)
		//将img矩阵转换成image图片
		backgroundImg, _ := img.ToImage()
		//获取整个img所占有的矩形区域Rectangle
		b := backgroundImg.Bounds()
		//NewRGBA 在矩形区域返回包含了 RGBA （红绿蓝透明）色值的一张图片
		m := image.NewRGBA(b)
		//绘制这张图片
		//参数1：目标图
		//参数2：矩形区域
		//参数3：已有图片
		//参数4：图片开始坐标，坐标系右下为坐标系的正方向
		//参数5：src作为前景图片存在
		draw.Draw(m, b, backgroundImg, image.ZP, draw.Src)
		//最大的人脸矩形区域（rect的面积）
		var maxArea = 0
		//最大的人脸的索引
		var maxIndex = 0
		//取出矩形区域的最大面积，和此最大面积的索引值
		for i := 0; i < len(rects); i++ {
			rectangle := rects[i]
			//获取x轴和y轴长度，获取其面积区域
			area := rectangle.Dx() * rectangle.Dy()
			if maxArea < area {
				area = maxArea
				maxIndex = i
			}
		}
		//rects大于0，表示人脸的数量大于0的
		if len(rects) > 0 {
			//获取判断到的最大的矩形，最大的那张脸
			r := rects[maxIndex]
			// 矩形左上角（向上偏移80个像素）
			startPoint := image.Pt(r.Min.X, r.Min.Y-80)
			// 矩形右下角（向右下偏移各自50，20个像素）
			endPoint := image.Pt(r.Max.X+50, r.Max.Y+20)

			// 由上述2个点重新构造矩形（矩形区域太小的话，狗头图片不足以覆盖整张脸）
			r = image.Rectangle{Min: startPoint, Max: endPoint}

			//将狗头图片，截取一块，按照x和y的大小进行图片的等比例拉伸，
			srcPng = imaging.Fill(srcPng, r.Size().X, r.Size().Y, imaging.Center, imaging.Lanczos)
			//将srcPng绘制在m上，也就是将截取后的狗头图片，放置在相机捕获的矩形区域内
			draw.Draw(m, r, srcPng, image.ZP, draw.Over)
		}
		//创建一个全新的矩阵，用于存储图片
		mat := gocv.NewMat()
		//将m转换成矩阵
		mat, err = gocv.ImageToMatRGB(m)
		if mat.Empty() {
			continue
		}
		//将矩阵展示在窗口
		window.IMShow(mat)

		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
