package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/epub-extractor/utils"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/extract", extractHandler)
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, nil)
}

func extractHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB

	// 从表单中获取上传的epub文件
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// 保存文件到临时文件夹中
	tempFile, err := os.CreateTemp("", "*.epub")
	if err != nil {
		log.Println(err)
		return
	}
	defer os.Remove(tempFile.Name())

	originalFileName := r.FormValue("filename")
	log.Println(originalFileName)
	_, err = io.Copy(tempFile, file)
	if err != nil {
		log.Println(err)
		return
	}

	// 调用提取函数提取epub中的文本和图片
	utils.LoadEpub(tempFile.Name(), originalFileName)

	// fmt.Fprintln(w, "Extraction complete.")
	http.ServeFile(w, r, "templates/success.html")
}
