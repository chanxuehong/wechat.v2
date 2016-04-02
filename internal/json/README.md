encoding/json 默认会对 html 标记 <, >, & 做编码转换(转换为 \u003c, \u003e, \u0026),  
但是狗日的腾讯不识别, 所以只能 hack encoding/json 库, 去掉这三个字符的转换.  

**对应 golang 1.6, 修改的文件有:**  
1. encode.go: 注释掉 794, 870 行的部分代码  
2. indent.go: 注释掉 21-29 行的全部代码  