package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Name    string
	Suffix  string
	Size    int64
	Type    string
	ModTime string
}

type Directory struct {
	Name        string
	Path        string
	Files       []File
	Directories []Directory
}

func GetDicByPath(path string, dic *Directory) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic("错误：" + err.Error())
	}
	for _, file := range files {
		if file.IsDir() {
			dicPath := filepath.Join(dic.Path, strings.TrimRight(file.Name(), "/"))
			dicInfo := InitDicInfo(dicPath)
			GetDicByPath(dicPath, &dicInfo)
			dic.Directories = append(dic.Directories, dicInfo)
		} else {
			InitFileInfo(file, dic)
		}
	}
}

func InitDicInfo(path string) Directory {
	info, err := os.Stat(path)
	if err != nil {
		panic("错误：" + err.Error())
	}
	var dic = Directory{
		Name:        info.Name(),
		Path:        path,
		Files:       make([]File, 0),
		Directories: make([]Directory, 0),
	}
	return dic
}
func InitFileInfo(file os.DirEntry, dic *Directory) {
	fileInfo, _ := file.Info()
	str := strings.Split(fileInfo.Name(), ".")
	suffix := ""
	if cap(str) > 1 {
		suffix = str[len(str)-1]
	}
	f := File{
		Name:    fileInfo.Name(),
		Suffix:  suffix,
		Size:    fileInfo.Size() / 1048576,
		ModTime: fileInfo.ModTime().Format("2006-01-02"),
	}
	dic.Files = append(dic.Files, f)
}

func (dic Directory) PrintDic() {
	fmt.Printf("文件夹:{ 名称:: %s 路径: %s}\n", dic.Name, dic.Path)
	for _, f := range dic.Files {
		fmt.Printf("-------- 文件:{ 【名称】: %s 【后缀】: %s 【大小】: %d MB 【修改时间】: %v }\n", f.Name, f.Suffix, f.Size, f.ModTime)
	}
	for _, subDir := range dic.Directories {
		subDir.PrintDic()
	}
}

func main() {
	var dirPath string
	flag.StringVar(&dirPath, "path", "E:/迅雷下载", "文件夹路径")
	flag.Parse()

	if dirPath == "" {
		panic("ERROR: 请输入文件夹路径")
	}
	dic := Directory{
		Name: dirPath,
		Path: dirPath,
	}
	GetDicByPath(dirPath, &dic)
	dic.PrintDic()
}
