package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

// 字段名必须要大写
type Config struct {
	FilePath string `conf:"file_path"`
	FileName string `conf:"file_name"`
	MaxSize  int64  `conf:"max_size"`
}

// 从conf文件中读取内容赋值给结构体指针
func parseConf(fileName string, result interface{}) (err error) {
	// 0 . result必须是一个指针
	t := reflect.TypeOf(result)
	v := reflect.ValueOf(result)
	if t.Kind() != reflect.Ptr {
		err = errors.New("result必须是一个指直针")
		return
	}
	// 判断是不是结构体
	if t.Elem().Kind() != reflect.Struct {
		err = errors.New("result必须是一个结构体指直针")
		return
	}

	// 1. 打开文件
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		//err = errors.New("打开配置文件失败")
		err = fmt.Errorf("打开配置文件%s失败", fileName)
		return
	}
	// 将读到的文件数据进行分割
	lineSlice := strings.Split(string(data), "\r\n")
	// 一行一行的解析
	for index, line := range lineSlice {
		line = strings.TrimSpace(line) // 去除字符串首尾的空白
		// 忽略空行和注释
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			// 是注释
			continue
		}
		// 解析是不是对的配置项, 判断是不是有等号
		equalIndex := strings.Index(line, "=")
		//fmt.Println("=====")
		//fmt.Println(equalIndex)
		//fmt.Println("======")
		if equalIndex == -1 {
			err = fmt.Errorf("第%d行语法错误", index+1)
			return
		}
		// 按照=分割每一行， 左边是key, 右边是value
		key := line[:equalIndex]
		key = strings.TrimSpace(key)

		value := line[equalIndex+1:]
		value = strings.TrimSpace(value)
		if len(key) == 0 {
			err = fmt.Errorf("第%d行语法错误,", index+1)
			return
		}
		// 利用反射给result赋值
		// 遍历结构体的每一个字段和key比较， 匹配上了就把value赋值
		for i := 0; i < t.Elem().NumField(); i++ {
			filed := t.Elem().Field(i)
			tag := filed.Tag.Get("conf")
			if key == tag {
				// 匹配上了就把value赋值
				// 拿到每个字段的类型
				filedType := filed.Type
				switch filedType.Kind() {
				case reflect.String:
					filedValue := v.Elem().FieldByName(filed.Name) // 根据字段名找到对应的字段值
					filedValue.SetString(value)                    // 将配置文件中读取的value设置给结构体

				case reflect.Int64:
					value64, _ := strconv.ParseInt(value, 10, 64)
					v.Elem().Field(i).SetInt(value64)
				}

			}
		}

	}

	return
}

func main() {

	// 2. 读取文件
	// 3. 读取每一行内容， 根据tag找结构体里面对应字段
	// 4. 找到了要赋值

	var c = &Config{} // 用来存储读取的数据
	err := parseConf("log.conf", c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
	fmt.Println(c.FileName)
}
