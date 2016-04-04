## 针对腾讯json的一些限制(bug?), 对 go1.6 encoding/json 库做一些hack

encoding/json 默认会对 html 标记 <, >, & 做转义(转义为 \u003c, \u003e, \u0026),  
但是狗日的腾讯不识别, 所以只能 hack encoding/json 库, 去掉这三个字符的转义.  

1. encode.go:  注释掉 794, 870 行的部分代码  
2. indent.go:  注释掉 21-29 行的全部代码  
  

狗日的腾讯返回的用户资料的昵称, 城市, 省份字段可能包含控制字符(我也不知道这些字段怎么会有控制字符, bug?),  
并且这些控制字符没有做 \uXXXX 转义, 导致 encoding/json 库解码失败, 再次 hack encoding/json 库!  

1. scanner.go: 注释掉 340-342 行的全部代码  

TODO: 貌似腾讯对所有的 \uXXXX 都不识别, encoding/json 默认对 控制字符, U+2028 和 U+2029 也做了转义, 这里并没有去掉, 等发现问题了再hack吧. 
