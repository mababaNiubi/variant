# variant

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[English](README.md)

Go 语言的通用动态类型容器使数据变成弱类型。**Variant** 是一个可辨识联合类型，支持八种运行时类型 —— Empty、Bool、Int、UInt、Float、String、List、Map，并提供 JSON/自定义二进制序列化、混合类型算术运算以及基于反射的 Go 原生类型解码。

专为数据传输、采集、遥测管道等场景设计，数据常以弱类型字符串的形式到达，需要高效的数值计算和存储。

## 安装

```bash
go get github.com/mababaNiubi/variant
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/mababaNiubi/variant"
)

func main() {
    // 智能构造器 —— 自动检测类型
    v := variant.New("3.14")
    fmt.Println(v.Type())        // TypeFloat64
    fmt.Println(v.AsFloat64())   // 3.14 <nil>

    // JSON 解析
    v2 := variant.New(`{"name":"Alice","age":30}`)
    name, _ := v2.MapGet("name")
    fmt.Println(name.AsString()) // "Alice"

    // 算术运算
    a := variant.NewInt(10)
    b := variant.NewFloat64(3.5)
    result, _ := a.Increase(b)
    fmt.Println(result.AsFloat64()) // 13.5 <nil>

    // JSON 序列化
    data, _ := v2.MarshalJSON()
    fmt.Println(string(data))    // {"age":30,"name":"Alice"}
}
```

## 构造器

```go
v := variant.New(x)          // 智能工厂 —— 接受任意 Go 值
v := variant.NewEmpty()      // TypeEmpty
v := variant.NewBool(true)   // TypeBool
v := variant.NewInt(42)      // TypeInt64
v := variant.NewInt64(42)    // TypeInt64
v := variant.NewUInt64(42)   // TypeUInt64
v := variant.NewValue(3.14)  // 自动判断整数还是浮点
v := variant.NewFloat64(3.14)// TypeFloat64
v := variant.NewString("hi") // TypeString
v := variant.NewValueList([]variant.Variant{...})   // TypeList
v := variant.NewValueMap(map[string]variant.Variant{...}) // TypeMap
```

`New()` 接受的类型: `bool`, `string`, `float64`, `float32`, `int`, `int8`–`int64`, `uint`–`uint64`, `[]byte`, `Variant`, `*Variant`, `[]Variant`, `map[string]Variant`，以及任意结构体（通过反射）。

对于字符串输入，`New()` 按以下优先级尝试解析：int → float → empty → JSON → 普通字符串。

## 类型转换

```go
b, err := v.AsBool()       // → bool
i, err := v.AsInt()        // → int
i64, err := v.AsInt64()    // → int64
u64, err := v.AsUInt64()   // → uint64
f32, err := v.AsFloat32()  // → float32
f64, err := v.AsFloat64()  // → float64
s := v.AsString()          // → string (永不失败)
s := v.String()            // → string (实现 fmt.Stringer 接口)
iface := v.AsInterface()   // → 递归转为 Go 原生类型
```

类型转换支持跨类型自动强制，并带有溢出保护。

## 相等性与比较

```go
v.IsEqual(other)                              // 深度比较
ok, err := v.CompareNumberBySymbol(r, ">=")   // 数值比较
ok := v.Comparable(r)                         // 数值比较，失败时降级为字典序
```

比较符号: `=`、`!=`、`>`、`<`、`>=`、`<=`。

## 算术运算

```go
result, err := a.Increase(b)  // a + b
result, err := a.Reduce(b)    // a - b
result, err := a.Multiple(b)  // a * b
result, err := a.Divide(b)    // a / b
result := v.Decimal(2)       // 保留 2 位小数
```

混合类型规则：

- Int + Float → Float
- String + 任意 → 字符串拼接
- List + List → 合并列表
- List + 标量 → 追加元素
- Map + Map → 合并映射（重复键覆盖）

## 容器操作

```go
// 多态方法（接受 int 索引或 string 键）
v.Add(value)               // 追加到列表或拼接到字符串
v.Set(indexOrKey, value)   // 按索引（列表）或键（映射）设置
val, ok := v.Get(idxOrKey) // 按索引或键获取
v.Remove(idxOrKey)         // 按索引或键删除

// 类型专属方法
val, ok := v.ListGet(0)
v.ListSet(0, value)
val, ok := v.MapGet("key")
v.MapSet("key", value)     // 空 Variant 上调用会初始化为 TypeMap

// 遍历
v.Range(func(key string, val Variant) bool { ... })
v.RangeByIndex(func(idx int, val Variant) bool { ... })
```

回调返回 `false` 可提前终止遍历。

## JSON 序列化

```go
// 序列化
data, err := v.MarshalJSON()
json.Marshal(v)  // 通过 json.Marshaler 接口

// 反序列化
var v variant.Variant
v.UnmarshalJSON([]byte(`{"key":"value"}`))

// 独立函数
v, err := variant.UnmarshalJSON([]byte(`[1,2,3]`))
```

## 二进制序列化

```go
// 带格式标记（适用于 WAL / 存储）
data, err := v.MarshalBinary()
v, n, err := variant.UnmarshalBinary(data)
ok := variant.IsBinaryFormat(data)

// 零分配批量编码 — 复用同一个 buffer
var buf []byte
for _, v := range manyVariants {
    buf = v.AppendBinary(buf)
}
```

二进制格式采用紧凑的自定义编码：
- 1 字节格式标记 (`0x01`) 用于区分 JSON
- 1 字节类型标签 + 固定宽度标量负载（所有数值类型 8 字节）
- 字符串/容器类型使用 `uint32` 长度前缀

零外部依赖 — 纯 Go 标准库实现。

## 解码 Go 结构体

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
v := variant.New(Person{Name: "Alice", Age: 30})
// v 为 TypeMap, 键为 "name" 和 "age" (来自 json tag)
```

使用反射实现。结构体字段优先使用 `json` tag 作为键名，无 tag 时使用 Go 字段名。

## 性能设计

- **内联数值存储**：整数、布尔值和 float64 通过 `unsafe.Pointer` 位转换存储在 `int64` 字段中 —— 常见数值路径无堆分配。
- **零分配二进制编码**：`AppendBinary` 直接写入调用者提供的 `[]byte` 缓冲区，无中间 buffer、无 encoder 对象、无第三方库开销。每个标量值只需一次 `append` 调用。
- **零分配二进制解码**：`UnmarshalBinary` 解析类型标签并直接通过 `encoding/binary` 读取固定宽度负载 —— 无 `interface{}` 装箱、无 decoder 分配。
- **延迟字符串解析**：数值字符串在被请求类型转换前不会执行解析。

## License

MIT — 详见 [LICENSE](LICENSE)。
