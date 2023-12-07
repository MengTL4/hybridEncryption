package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PaddingFunc 是用于填充数据的函数类型，接收一个 []byte 和一个块大小参数，返回填充后的 []byte。
type PaddingFunc func([]byte, int) []byte

// UnpaddingFunc 是用于去除填充数据的函数类型，接收一个 []byte，返回去除填充后的 []byte。
type UnpaddingFunc func([]byte) []byte

// ReadStringHex 从文件中读取十六进制字符串并返回。
func ReadStringHex(filename string) string {
	f, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("读取", filename, "时发生错误", err)
	}
	return string(f)
}

// WriteStringHex 将十六进制字符串写入文件。
func WriteStringHex(filename string, msg string) {
	err := os.WriteFile(filename, []byte(msg), 0666)
	if err != nil {
		fmt.Println("写入", filename, "时发生错误", err)
	}
}

// ReadBytesHex 从文件中读取十六进制数据并返回字节数组。
func ReadBytesHex(filename string) []byte {
	tmp, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("读取", filename, "时发生错误", err)
	}
	if len(tmp)%2 == 1 {
		tmp = append(tmp[:len(tmp)-1], byte('0'), tmp[len(tmp)-1])
	}
	hexStringReader := strings.NewReader(string(tmp))
	tmp2 := make([]byte, len(tmp)/2)
	_, err = fmt.Fscanf(hexStringReader, "%x", &tmp2)
	return tmp2
}

// WriteBytesHex 将字节数组以十六进制格式写入文件。
func WriteBytesHex(filename string, msg []byte) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("打开文件错误 =", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%x", msg)
	if err != nil {
		return
	}
	return
}

// DumpWords 打印给定的无符号整数数组，每行显示四个元素。
func DumpWords(note string, in []uint32) {
	fmt.Printf("%s", note)
	for i, v := range in {
		if i%4 == 0 {
			fmt.Printf("\nword[%02d]: %.8x ", i/4, v)
		} else {
			fmt.Printf("%.8x ", v)
		}
	}
	fmt.Println("\n")
}

// DumpBytes 打印给定的字节数组，每行显示十六个元素。
func DumpBytes(note string, in []byte) {
	fmt.Printf("%s", note)
	for i, v := range in {
		if i%16 == 0 {
			fmt.Printf("\nblock[%d]: %02x", i/16, v)
		} else {
			if i%4 == 0 {
				fmt.Printf(" %02x", v)
			} else {
				fmt.Printf("%02x", v)
			}
		}
	}
	fmt.Println("\n")
}

// ZeroPadding 使用零填充方式对字节数组进行填充，以满足给定的块大小。
func ZeroPadding(in []byte, blockLen int) []byte {
	tmp := make([]byte, len(in))
	copy(tmp, in)

	remainder := len(tmp) % blockLen
	for i := 0; i < blockLen-remainder; i++ {
		tmp = append(tmp, 0x00)
	}
	return tmp
}

// ZeroUnpadding 从字节数组中移除零填充。
func ZeroUnpadding(in []byte) []byte {
	for in[len(in)-1] == 0x00 {
		in = in[:len(in)-1]
	}
	tmp := make([]byte, len(in))
	copy(tmp, in)
	return tmp
}

// PKCS7Padding 使用PKCS7填充方式对字节数组进行填充，以满足给定的块大小。
func PKCS7Padding(in []byte, blockLen int) []byte {
	tmp := make([]byte, len(in))
	copy(tmp, in)

	rmd := len(tmp) % blockLen
	for i := 0; i < blockLen-rmd; i++ {
		tmp = append(tmp, byte(blockLen-rmd))
	}
	return tmp
}

// PKCS7Unpadding 从字节数组中移除PKCS7填充。
func PKCS7Unpadding(in []byte) []byte {
	last := int(in[len(in)-1])
	tmp := make([]byte, len(in)-last)
	copy(tmp, in[:len(in)-last])
	return tmp
}

// 定义一个创建文件目录的方法
func Mkdir(basePath string) string {
	//	1.获取当前时间,并且格式化时间
	folderName := time.Now().Format("2006/01/02")
	folderPath := filepath.Join(basePath, folderName)
	//使用mkdirall会创建多层级目录
	os.MkdirAll(folderPath, os.ModePerm)
	return folderPath
}
