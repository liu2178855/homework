package main

import (
	"os"
	"io"
	"fmt"
	"flag"
	"bufio"
	"strings"
	"hash/crc32"
)

var fileName string
var topn     int
var chunk    int

func BufWriterFile(ff string, mes string) {
	fp, err := os.OpenFile(ff, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0751)
	if err != nil {
		fmt.Printf("open file %s failed: %v\n", ff, err)
		panic(err)
	}
	defer fp.Close()
	writer := bufio.NewWriter(fp)
	_, err = writer.Write([]byte(mes))
	if err != nil {
		fmt.Println("write file %s failed: %v\n", ff, err)
	}
	writer.Flush()
}

func getFileDictArr(ff string) []dict {
	fp, err := os.OpenFile(ff, os.O_RDONLY, 0751)
        if err != nil {
                fmt.Printf("open file %s failed: %v\n", fileName, err)
                panic(err)
        }
	defer fp.Close()
	reader := bufio.NewReader(fp)
        var rstr string
        trie := NewTrie()
        for {
                rstr, err = reader.ReadString('\n')
                if err == nil {
                        str := strings.TrimSuffix(rstr, "\n")
                        trie.Insert(str)				//将子文件内容写入字典树，区分度较高的情况可以写入map，实现简单，速度快
                } else if err == io.EOF {
                        break
                } else {
                        fmt.Printf("read file %s failed: %v\n", fileName, err)
                        panic(err)
                }
        }
	arr := trie.GetDictArr()
	return arr
}

func compare(a []dict, b []dict) []dict {
	lenth1 := len(a)
	lenth2 := len(b)
	var i, j, cnt int
	var res []dict
	for {
		if cnt >= topn || (i >= lenth1 && j >= lenth2) {
			break
		}
		if i < lenth1 && j < lenth2 {
			if a[i].num > b[j].num {
				res = append(res, a[i])
				i ++
			} else {
				res = append(res, b[j])
				j ++
			}
		} else if i >= lenth1 {
			res = append(res, b[j])
			j ++
		} else if j >= lenth2 {
			res = append(res, a[i])
                        i ++
		}
		cnt++
	}
	return res
}

func Exist(ff string) bool {
    _, err := os.Stat(ff)
    return err == nil || os.IsExist(err)
}

func main() {
	flag.StringVar(&fileName, "file", "", "the file name you choose")
	flag.IntVar(&topn, "topn", 10, "output the top n numbers")
	flag.IntVar(&chunk, "chunk", 1024, "split file number")
	flag.Parse()
	fp, err := os.OpenFile(fileName, os.O_RDONLY, 0751)
	if err != nil {
		fmt.Printf("open file %s failed: %v\n", fileName, err)
		panic(err)
	}
	reader := bufio.NewReader(fp)
	var rstr string
	for {							//按hash分发源文件的字符串到子文件
		rstr, err = reader.ReadString('\n')
		if err == nil {
			str := strings.TrimSuffix(rstr, "\n")
			hashcode := crc32.ChecksumIEEE([]byte(str))
			fileStr := fileName + "_" + fmt.Sprintf("%04d", hashcode%uint32(chunk))
			BufWriterFile(fileStr, rstr)
		} else if err == io.EOF {
			break
		} else {
			fmt.Printf("read file %s failed: %v\n", fileName, err)
			panic(err)
		}
	}
	fp.Close()
	fmt.Println("dispatch finishied")
	var maxArr []dict					//maxArr存当前topn的字符串及其出现次数
	for i := 0; i < chunk; i++ {
		fileStr := fileName + "_" + fmt.Sprintf("%04d", i)
		if !Exist(fileStr) {
			continue
		}
		tmp := getFileDictArr(fileStr)
		maxArr = compare(maxArr, tmp)
		fmt.Printf("finish %s\n", fileStr)
	}
	for _, v := range maxArr {
		fmt.Println(v)
	}
}
