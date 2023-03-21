package utils

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

func extractText(epubName string, r io.Reader) {
	epubName = sanitizeFileName(epubName)
	filepath := filepath.Join("text", epubName+".txt")
	doc := etree.NewDocument()
	_, err := doc.ReadFrom(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 提取所有p标签的文本
	var text string
	for _, p := range doc.FindElements("//p") {
		text += extractTextFromElement(p)
	}

	// 保存到text.txt文件中
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(text)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.Flush()
}

// 保存图片文件到pictures目录中
func saveImage(name string, epubName string, r io.Reader) {
	// 只保留图片文件名，去掉前缀
	name = filepath.Base(name)

	// 创建保存图片的文件
	f, err := os.Create(filepath.Join("pictures", epubName, name))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// 复制数据到文件中
	_, err = io.Copy(f, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func LoadEpub(epubFile string, epubName string) {
	// 打开epub文件
	r, err := zip.OpenReader(epubFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Close()

	// 创建保存图片的目录
	err = os.MkdirAll("pictures", os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.MkdirAll("pictures/"+epubName, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建保存文本的目录
	err = os.MkdirAll("text", os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	//删除旧的文本
	err = deleteFileIfExists("text/" + sanitizeFileName(epubName) + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 遍历epub文件中的所有文件
	for _, f := range r.File {
		// 只处理html文件和图片文件
		if strings.HasSuffix(f.Name, ".html") || strings.HasSuffix(f.Name, ".xhtml") || strings.HasPrefix(f.Name, "OEBPS/Images/") {
			// 打开文件
			rc, err := f.Open()
			if err != nil {
				fmt.Println(err)
				continue
			}
			defer rc.Close()

			// 提取文本
			if strings.HasSuffix(f.Name, ".xhtml") || strings.HasSuffix(f.Name, ".html") {
				extractText(epubName, rc)
			}

			// 提取图片
			if strings.HasPrefix(f.Name, "OEBPS/Images/") {
				saveImage(f.Name, epubName, rc)
			}
		}
	}

	fmt.Println("Extraction complete.")
}

// 替换不允许作为文件名的字符
func sanitizeFileName(name string) string {
	r := strings.NewReplacer("<", "_",
		">", "_",
		":", "_",
		"\"", "_",
		"/", "_",
		"\\", "_",
		"|", "_",
		"?", "_",
		"*", "_")
	return r.Replace(name)
}

func deleteFileIfExists(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 文件不存在，不需要删除
		return nil
	}

	// 文件存在，删除文件
	return os.Remove(filename)
}

func extractTextFromElement(elem *etree.Element) string {
	text := elem.Text()
	for _, child := range elem.ChildElements() {
		text += extractTextFromElement(child)
	}
	text += elem.Tail()
	return text
}
