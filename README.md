# Tableau
A powerful configuration conversion tool based on Protobuf(proto3).

## Features
- Convert **xlsx** to **JSON**, JSON is the first-class citizen of exporting targets.
- Use **protobuf** as the IDL(Interface Description Language) to define the structure of **xlsx**.
- Use **golang** to develop the conversion engine.
- Support multiple programming languages, thanks to **protobuf**.

## Concepts
- Importer: xlsx importer.
- IR: Intermediate Representation.
- Filter: filter the IR.
- Exporter: JSON exporter, protobin exporter, prototext exporter, xml exporter, sqlite3 exporter, and so on.
- ProtoConf: a configuration metadata format based on protobuf.

## Types
- Scalar
- Message(struct)
- List
- Map(unordered)
- Timestamp
- Duration

## TODO

### protoc plugins
- [ ] Golang
- [ ] C++
- [ ] C#/.NET
- [ ] Python
- [ ] Lua
- [ ] Javascript/Typescript/Node
- [ ] Java

### Metadata
- [ ] metatable: a message to describe the worksheet's metadata.
- [ ] metafield: a message to describe the caption's metadata.
- [x] captrow: caption row, the exact row number of captions at worksheet. **Newline** in caption is allowed for more readability, and will be trimmed in conversion. 
- [ ] descrow: description row, the exact row number of descriptions at worksheet.
- [x] datarow: data row, the start row number of data.

[Newline](https://www.wikiwand.com/en/Newline)(line break) in major operating systems:

| OS                  | Abbreviation | Escape sequence |
| ------------------- | ------------ | --------------- |
| Unix (linux, OS X)  | LF           | `\n`            |
| Microsoft Windows   | CRLF         | `\r\n`          |
| classic Mac OS/OS X | CR           | `\r`            |

> **LF**: Line Feed, **CR**: Carriage Return.
>
> [Mac OS X](https://www.oreilly.com/library/view/mac-os-x/0596004605/ch01s06.html)

### Generator
- [x] generate xlsx template by proto: **proto -> xlsx**
- [x] generated xlsx template caption row with style: font bold, backgroud color and so on. See [XLSX Styles](https://github.com/tealeg/xlsx/blob/master/tutorial/tutorial.adoc#styles)
- [ ] generate proto by xlsx template: **xlsx -> proto**

### Conversion
- [x] xlsx -> JSON(common format and human readable)
- [x] xlsx -> protobin(small size)
- [x] xlsx -> prototext(human debugging)
- [ ] JSON -> xlsx
- [ ] protobin -> xlsx
- [ ] prototext -> xlsx

### Pretty Print
- [x] Multiline: every textual element on a new line
- [x] Indent: 4 space characters
- [x] JSON support
- [x] prototext support

### EmitUnpopulated
- [x] JSON: `EmitUnpopulated` specifies whether to emit unpopulated fields.

### Scalar Types
- [x] interger: int32, uint32, int64 and uint64
- [x] float: float and double
- [x] bool
- [x] string
- [x] bytes

### Enumerations
- [ ] enum: The name of the enum value as specified in proto is used. Parsers accept both enum names and integer values. 
- [ ] validate the enum value.

### Composite Types
- [x] message: horizontal(row direction) layout, fields located in cells.
- [x] message: simple in-cell message, each field must be **scalar** type. It is a comma-separated list of fields. E.g.: `1,test,3.0`. List's size need not to be equal to fields' size, as fields will be filled in order. Fields not configured will be filled default values due to its scalar type.
- [x] list: horizontal(row direction) layout, and is list's default layout.
- [x] list: vertical(column direction) layout.
- [x] list: simple in-cell list, element must be **scalar** type. It is a comma-separated list of elements. E.g.: `1,2,3`. 
- [x] list: scalable or dynamic list size.
- [x] list: smart recognition of empty element at any position.
- [x] map: horizontal(row direction) layout.
- [x] map: vertical(column direction) layout, and is map's default layout.
- [x] map: unordered map or hash map.
- [ ] map: ordered map.
- [x] map: simple in-cell map, both key and value must be **scalar** type. It is a comma-separated list of key=value pairs. E.g.: `1:10,2:20,3:30`. 
- [x] map: scalable or dynamic map size.
- [x] map: smart recognition of empty value at any position.
- [x] nesting: unlimited nesting of message, list, and map.

### Default Values
Each scalar type's default value is same as protobuf.

- [x] interger: `0` 
- [x] float: `0.0` 
- [x] bool: `false`
- [x] string: `""`
- [x] bytes: `""`
- [x] in-cell message: each field's default value is same as protobuf 
- [x] in-cell list: element's default value is same as protobuf 
- [x] in-cell map: both key and value's default value are same as protobuf 
- [x] message: all fields have default values

### Empty
- [x] scalar: default value same as protobuf.
- [x] message: empty message will not be spawned if all fields are empty.
- [x] list: empty list will not be spawned if list's size is 0.
- [x] list: empty message will not be appended if list's element(message type) is empty.
- [x] map: empty map will not be spawned if map's size is 0.
- [x] map: empty message will not be inserted if map's value(message type) is empty.
- [x] nesting: recursively empty.

### Merge
- [ ] merge multiple workbooks
- [ ] merge multiple worksheets

### Datetime
> [Understanding about RFC 3339 for Datetime and Timezone Formatting in Software Engineering](https://medium.com/easyread/understanding-about-rfc-3339-for-datetime-formatting-in-software-engineering-940aa5d5f68a)
> ```
> # This is acceptable in ISO 8601 and RFC 3339 (with T)
> 2019-10-12T07:20:50.52Z
> # This is only accepted in RFC 3339 (without T)
> 2019-10-12 07:20:50.52Z
> ```
> - "Z" stands for **Zero timezone** or **Zulu timezone** `UTC+0`, and equal to `+00:00` in the RFC 3339.
> - **RFC 3339** follows the **ISO 8601** DateTime format. The only difference is RFC allows us to replace "T" with "space".

Use [RFC 3339](https://tools.ietf.org/html/rfc3339) , which is following [ISO 8601](https://www.wikiwand.com/en/ISO_8601).

- [x] Timestamp: based on `google.protobuf.Timestamp`, see [JSON mapping](https://developers.google.com/protocol-buffers/docs/proto3#json)
- [x] Timezone: see [ParseInLocation](https://golang.org/pkg/time/#ParseInLocation)
- [ ] DST: Daylight Savings Time. *There is no plan to handle this boring stuff*.
- [x] Datetime: excel format: `yyyy-MM-dd HH:mm:ss`, e.g.: `2020-01-01 05:10:00`
- [ ] Date: excel format: `yyyy-MM-dd`, e.g.: `2020-01-01`
- [ ] Time: excel format: `HH:mm:ss`, e.g.: `05:10:00`
- [x] Duration: based on`google.protobuf.Duration` , see [JSON mapping](https://developers.google.com/protocol-buffers/docs/proto3#json)
- [x] Duration: excel format: `form "72h3m0.5s"`, see [golang duration string form](https://golang.org/pkg/time/#Duration.String)
  
### Transpose
- [x] Interchange the rows and columns of a worksheet.

### Validation
- [ ] Min
- [ ] Max
- [ ] Range
- [ ] Options: e.g.: enum type
- [ ] Foreign key

### Error Message
- [ ] Report clear and precise error messages when converter failed, please refer to the programming language compiler
- [ ] Use golang template to define error message template
- [ ] Multiple languages support, focused on English and Simplified Chinese

### Performace
- [ ] Stress test
- [ ] Each goroutine process one worksheet
- [ ] Mutiple process model

### Optimization
- [ ] Error: [https://github.com/pkg/errors](https://github.com/pkg/errors)
- [ ] Log: [https://github.com/uber-go/zap](https://github.com/uber-go/zap)
