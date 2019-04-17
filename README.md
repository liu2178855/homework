## 安装：

	go build main.go trie.go


## 使用：

	./main --file=benchmark --topn=100 --chunk=1024

* file参数表示源url文件
* topn参数表示输出重复出现前topn的字符串及次数
* chunk参数表示通过每个字符串的hash值路由到多少个子文件（根据源文件大小以及内存大小限制来确定这个参数）

## 执行流程

1.源文件通过crc32算出hash值然后对chunk取模写入到相应的子文件。所以相同的字符串肯定会被写进同一个子文件当中
2.遍历每个子文件，把子文件字符串写到字典树中，然后对字符串及出现次数做排序，最终取得topn的字符串（这里考虑使用字典树是在前缀大部分相同的情况下更节省些内存，不过也可以用map存储，实现更简单，而且在字符串区分度很大的情况下占用内存更小，速度更快）
