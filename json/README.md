## 针对腾讯json的一些限制(bug?), 对 go1.6 encoding/json 库做一些hack

#### encoding/json 默认会对 html 标记 <, >, & 做转义(转义为 \u003c, \u003e, \u0026), 但是腾讯不识别, 所以只能 hack encoding/json 库, 去掉这三个字符的转义.

**TODO: 貌似腾讯对所有的 \uXXXX 都不识别, encoding/json 默认对 控制字符(0x00-0x1F), U+2028 和 U+2029 也做了转义, 这里并没有去掉, 等发现问题了再hack吧**

搜索关键字: '<', '>', '&', '\u003c', '\u003e', '\u0026', '\x3c', '\x3e', '\x26', 0x3c, 0x3e, 0x26, \u
```
1. encode.go:      注释掉 794 行的部分代码
2. encode.go:      注释掉 870 行的部分代码
3. indent.go:      注释掉 21-29 行的全部代码
4. encode_test.go: 重写 TestMarshalerEscaping
5. decode_test.go: 重写 TestEscape
```

#### 对于标准的json格式, 数字和布尔型的字段值不应该被双引号包含的, 比如 {"x": 123456, "y": true}, 但是腾讯某些api返回的json的整数被双引号包含了, hack encoding/json 库, 自适应这些双引号.

搜索关键字: "case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:"
```
1. hack.go:      新增 unmarshalStringToNumberOrBoolean 函数
2. hack_test.go: 新增 TestUnmarshalStringToNumberOrBoolean 函数
3. decode.go:    修改 819 行代码
```

#### 腾讯返回的用户资料的昵称, 城市, 省份字段可能包含控制字符(0x00-0x1F, 我也不知道这些字段怎么会有控制字符, bug?), 并且这些控制字符没有做 \uXXXX 转义, 导致 encoding/json 库解码失败, 再次 hack encoding/json 库.

搜索关键字: control, 31, 32, 0x20, ' ', '\u0020', '\x20'
```
1. decode.go:    注释掉 1067 行的部分代码
2. decode.go:    注释掉 1150 行的部分代码
3. scanner.go:   注释掉 340-342 行的全部代码
4. hack_test.go: 新增 TestUnmarshalControlCharacter 函数
```

#### 因为sdk对于json的编码和解码都是用了hacked json库, 对于encoding/json.Number类型就不能正常处理, hack it

```
1. hack.go:   新增 var jsonNumberType = reflect.TypeOf(json.Number(""))
2. encode.go: 修改 531 行代码
3. decode.go: 修改 853 行代码
```
