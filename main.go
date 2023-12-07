package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"hybrid_encryption/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	//encryptFile()
	r := gin.Default()
	// Create a route for the web page.
	r.LoadHTMLGlob("templates/*")
	r.Static("/upload", "./upload")
	r.Static("/encrypt_upload", "./encrypt_upload")
	r.Static("/decrypt_upload", "./decrypt_upload")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Handle encryption request
	r.POST("/encrypt", func(c *gin.Context) {
		plaintext := c.PostForm("plaintext")
		if plaintext != "" {
			key := []byte("SpP5j4DACmh5uEWR")
			ciphertext, roundState, roundState2, roundState3, roundState4 := HybridEncrypt(key, plaintext)
			// 将 round 数据转换为 JSON 格式
			roundJSON := make(map[int]string)
			for k, v := range roundState {
				roundJSON[k] = fmt.Sprintf("%v", v)
			}
			roundJSON2 := make(map[int]string)
			for k, v := range roundState2 {
				roundJSON2[k] = fmt.Sprintf("%v", v)
			}
			roundJSON3 := make(map[int]string)
			for k, v := range roundState3 {
				roundJSON3[k] = fmt.Sprintf("%v", v)
			}
			roundJSON4 := make(map[int]string)
			for k, v := range roundState4 {
				roundJSON4[k] = fmt.Sprintf("%v", v)
			}
			// 将 roundJSON 转换为 JSON 格式
			roundJSONString, err := json.Marshal(roundJSON)
			if err != nil {
				c.String(http.StatusInternalServerError, "JSON 编码失败")
				return
			}
			roundJSONString2, err := json.Marshal(roundJSON2)
			if err != nil {
				c.String(http.StatusInternalServerError, "JSON 编码失败")
				return
			}
			roundJSONString3, err := json.Marshal(roundJSON3)
			if err != nil {
				c.String(http.StatusInternalServerError, "JSON 编码失败")
				return
			}

			roundJSONString4, err := json.Marshal(roundJSON4)
			if err != nil {
				c.String(http.StatusInternalServerError, "JSON 编码失败")
				return
			}
			// 构建返回的 JSON 数据
			responseJSON := fmt.Sprintf(`{"ciphertext": %s, "roundState": %s,"roundState2":%s,"roundState3":%s,"roundState4":%s}`, string(ciphertext), roundJSONString, roundJSONString2, roundJSONString3, roundJSONString4)
			c.String(http.StatusOK, responseJSON)
		} else {
			c.String(http.StatusOK, "请输入明文！")
		}
	})

	// Handle decryption request
	r.POST("/decrypt", func(c *gin.Context) {
		ciphertext := c.PostForm("ciphertext")
		decryptPlaintext := HybridDecrypt([]byte(ciphertext))
		c.String(http.StatusOK, decryptPlaintext)
	})
	r.POST("/detailFile", func(c *gin.Context) {
		if form, err := c.MultipartForm(); err == nil {
			//1.获取文件
			files := form.File["file"]
			//2.循环全部的文件
			operation := c.PostForm("operation")
			if operation == "encrypt" {
				var filePathString, encryptFilePath string
				for _, file := range files {
					// 3.根据时间鹾生成文件名
					fileNameInt := time.Now().Unix()
					fileNameStr := strconv.FormatInt(fileNameInt, 10)
					//4.新的文件名(如果是同时上传多张图片的时候就会同名，因此这里使用时间鹾加文件名方式)
					fileName := fileNameStr + file.Filename
					//5.保存上传文件
					filePathString = filepath.Join(utils.Mkdir("upload"), "/", fileName)
					c.SaveUploadedFile(file, filePathString)
					encryptFilePath = encryptFile(filePathString, fileName)
				}
				c.JSON(http.StatusOK, gin.H{
					"code":    0,
					"message": "加密成功",
					"path":    encryptFilePath,
				})
			} else {
				var filePathString, decryptFilePath string
				for _, file := range files {
					// 3.根据时间鹾生成文件名
					fileNameInt := time.Now().Unix()
					fileNameStr := strconv.FormatInt(fileNameInt, 10)
					//4.新的文件名(如果是同时上传多张图片的时候就会同名，因此这里使用时间鹾加文件名方式)
					fileName := fileNameStr + file.Filename
					//5.保存上传文件
					filePathString = filepath.Join(utils.Mkdir("upload"), "/", fileName)
					c.SaveUploadedFile(file, filePathString)
					decryptFilePath = decryptFile(filePathString, fileName)
				}
				c.JSON(http.StatusOK, gin.H{
					"code":    0,
					"message": "解密成功",
					"path":    decryptFilePath,
				})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "获取数据失败",
			})
		}
	})
	r.Run(":8000")
}
