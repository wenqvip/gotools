package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s + [filename]\n", os.Args[0])
		return
	}

	allFiles, err := getAllFiles()
	if err != nil {
		fmt.Printf("some thing wrong: %s\n", err)
		return
	}

	files := os.Args[1:]
	for _, filePattern := range files {
		fileList := getFileList(filePattern, allFiles)
		fmt.Printf("Pattern: %s, Total: %d\n", filePattern, len(fileList))
		for _, file := range fileList {
			fText, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("ioutil.ReadFile %s failed: %s\n", file, err)
				continue
			}

			charCode, err := detectCode(fText)
			if err != nil {
				fmt.Printf("detectCode failed: %s\n", err)
				continue
			}

			if charCode == "GB-18030" {
				newFile, err := os.OpenFile(file, os.O_RDWR, 0666)
				if err != nil {
					fmt.Printf("OpenFile %s failed: %s\n", file, err)
					newFile.Close()
					continue
				}

				// github.com/saintfish/chardet 只检测 GB-18030
				// golang.org/x/net/html/charset 只能用gbk
				newContent, err := convertToUtf8(fText, "gbk")
				if err != nil {
					fmt.Printf("convertToUtf8 failed: %s\n", err)
					newFile.Close()
					continue
				}

				prefix := []byte{0xEF, 0xBB, 0xBF}
				newContent = append(prefix, newContent...)
				_, err = newFile.Write(newContent)
				if err != nil {
					fmt.Printf("newFile.Write failed: %s\n", err)
					newFile.Close()
					continue
				}

				newFile.Close()
				fmt.Printf("%s convert from %s to UTF-8 BOM success!\n", file, charCode)
			} else if charCode == "UTF-8" {
				if len(fText) >= 3 && fText[0] == 0xEF && fText[1] == 0xBB && fText[2] == 0xBF {
					continue
				}
				newFile, err := os.OpenFile(file, os.O_RDWR, 0666)
				if err != nil {
					fmt.Printf("OpenFile %s failed: %s\n", file, err)
					newFile.Close()
					continue
				}

				prefix := []byte{0xEF, 0xBB, 0xBF}
				fText = append(prefix, fText...)
				_, err = newFile.Write(fText)

				if err != nil {
					fmt.Printf("newFile.Write failed: %s\n", err)
					newFile.Close()
					continue
				}

				newFile.Close()
				fmt.Printf("%s convert from %s to UTF-8 BOM success!\n", file, charCode)
			} else if strings.HasPrefix(charCode, "ISO-8859") {
				fmt.Printf("%s coded %s, do nothing!\n", file, charCode)
			} else {
				fmt.Printf("%s coded %s, do nothing!\n", file, charCode)
			}
		}
	}
}

func convertToUtf8(src []byte, encode string) ([]byte, error) {
	byteReader := bytes.NewReader(src)
	reader, err := charset.NewReaderLabel(encode, byteReader)
	if err != nil {
		fmt.Printf("charset.NewReaderLabel failed : %s\n", err)
		return nil, err
	}

	dst, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("ioutil.ReadAll failed : %s\n", err)
		return nil, err
	}
	return dst, nil
}

func detectCode(src []byte) (string, error) {
	detector := chardet.NewTextDetector()
	var result *chardet.Result
	result, err := detector.DetectBest(src)
	if err != nil {
		fmt.Printf("detector.DetectBest failed: %s\n", err)
		return "", err
	}

	fmt.Printf("charset: %s, language: %s, confidence: %d\n",
		result.Charset, result.Language, result.Confidence)

	return result.Charset, nil
}

func getFileList(pattern string, fileList []string) []string {
	var res = make([]string, 0, 10)
	for _, file := range fileList {
		_, filename := filepath.Split(file)
		if match, _ := filepath.Match(pattern, filename); match {
			res = append(res, file)
		}
	}
	return res
}

func getAllFiles() ([]string, error) {
	var allFiles = make([]string, 0, 100)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			allFiles = append(allFiles, path)
		}
		if err != nil {
			fmt.Printf("Walk err: %s\n", err)
			return err
		}
		return nil
	})

	return allFiles, err
}
