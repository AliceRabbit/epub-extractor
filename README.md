# EPUB提取工具
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

这是一个使用Go实现的EPUB提取工具，它可以从EPUB文件中提取文本和图片，并将它们分别保存在项目目录下的text和pictures文件夹中。此外，它还提供了一个简单的网页前端，可以让用户上传EPUB文件并使用提取工具进行提取。

## 使用方法
要使用此工具，请确保已经正确安装go。

克隆此仓库：
```bash
git clone https://github.com/AliceRabbit/epub-extractor.git
```
进入仓库目录：

```bash
cd epub-extractor
```
构建可执行文件：

```go
go build
```
运行可执行文件（Window下双击生成的exe文件）：

```bash
./epub-extractor
```
打开Web浏览器，访问http://localhost:8080，上传EPUB文件并进行提取。

## 注意事项
目前，此工具仅能提取.epub格式的文件，其他格式的文件将被忽略。

此工具使用了go-epub库来解析EPUB文件，请确保将此库安装到系统中。

如果您使用的是Windows系统，请将可执行文件重命名为epub-extractor.exe，并在终端中使用此名称运行它。

如果您需要将此工具部署到生产环境中，请注意加强安全措施，例如使用HTTPS、限制文件上传大小等。
当前版本没有提供多线程支持，任何并行操作都可能造成不确定的意外后果。