# golang脚本工具集合 ![Go](https://github.com/wenqvip/gotools/workflows/Go/badge.svg)

## 编译
```sh
make all
```

## gbk2utf
自动将当前目录下所有GBK编码的文件转换为UTF8编码，主要用来转换源码文件。由于VS对不带BOM的UTF8源文件不能正常识别，转化后的文件都是带BOM的

```sh
# usage
gbk2utf *.h *.cpp *.cc
```