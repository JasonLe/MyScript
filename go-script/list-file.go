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
	ModTime string
}

type Directory struct {
	Name        string
	Path        string
	Files       []File
	Directories []Directory
}

func GetDicByPath(path string, dic *Directory) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("读取目录错误：%w", err)
	}
	for _, file := range files {
		if file.IsDir() {
			dicPath := filepath.Join(dic.Path, strings.TrimRight(file.Name(), "/"))
			dicInfo, err := InitDicInfo(dicPath)
			if err != nil {
				return fmt.Errorf("初始化目录信息错误：%w", err)
			}
			GetDicByPath(dicPath, &dicInfo)
			dic.Directories = append(dic.Directories, dicInfo)
		} else {
			InitFileInfo(file, dic)
		}
	}
	return nil
}

func InitDicInfo(path string) (Directory, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Directory{}, fmt.Errorf("获取目录信息错误：%w", err)
	}
	dic := Directory{
		Name:        info.Name(),
		Path:        path,
		Files:       []File{},
		Directories: []Directory{},
	}
	return dic, nil
}

func InitFileInfo(file os.DirEntry, dic *Directory) {
	fileInfo, err := file.Info()
	if err != nil {
		fmt.Printf("获取文件信息错误：%v\n", err)
		return
	}
	suffix := filepath.Ext(fileInfo.Name())
	f := File{
		Name:    fileInfo.Name(),
		Suffix:  suffix,
		Size:    fileInfo.Size() / 1048576,
		ModTime: fileInfo.ModTime().Format("2006-01-02"),
	}
	dic.Files = append(dic.Files, f)
}

func (dic Directory) PrintDic() {
	fmt.Printf("【名称】: %s 【路径】: %s \n", dic.Name, dic.Path)
	for _, f := range dic.Files {
		fmt.Printf("-------- 【名称】: %s 【后缀】: %s 【大小】: %d MB 【修改时间】: %v \n", f.Name, f.Suffix, f.Size, f.ModTime)
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
		fmt.Println("ERROR: 请输入文件夹路径")
		return
	}
	dic := Directory{
		Name: dirPath,
		Path: dirPath,
	}
	err := GetDicByPath(dirPath, &dic)
	if err != nil {
		fmt.Println(err)
		return
	}
	dic.PrintDic()
}
