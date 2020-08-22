# Tableau
A powerful configuration conversion tool based on protobuf.

## Features
- Convert **xlsx** to **JSON**, JSON is the first-class citizen of exporting targets
- Comments in JSON: add one more comment key-value pair, and the comment key is prefixed with "#"
- Use **protobuf** as the IDL(Interface Description Language) to define the structure of **xlsx**
- Use **golang** to develop the conversion engine
- Multiple languages support, thanks to **protobuf**

## Concept
- Importer: xlsx importer
- IR: Intermediate Representation, in-memory object, JSON or protobin.
- Filter: filter the IR.
- Exporter: JSON exporter, protobin exporter, prototext exporter, xml exporter, sqlite3 exporter, and so on.

## TODO

### Testing
- [ ] Golang
- [ ] C++
- [ ] C#/.NET
- [ ] Python
- [ ] Lua
- [ ] Java
- [ ] Javascript/Typescript/Node

### Metadata
- [ ] metatable: a message to describe the worksheet's metadata.
- [ ] metafield: a message to describe the caption's metadata.
- [x] captrow: caption row, the exact row number of captions at worksheet.
- [ ] descrow: description row, the exact row number of descriptions at wooksheet.
- [x] datarow: data row, the start row number of data.

### Generator
- [ ] generate xlsx template by proto: **proto -> xlsx template**
- [ ] generate proto by xlsx template: **proto -> xlsx template**

### Conversion
- [x] xlsx -> JSON(common format and human readable)
- [x] xlsx -> protobin(small size)
- [x] xlsx -> prototext(human debugging)
- [ ] JSON -> xlsx
- [ ] protobin -> xlsx
- [ ] prototext -> xlsx

### Pretty Print
- [x] JSON
- [x] prototext

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
- [x] message: simple in-cell message, each field must be **scalar** type.
- [x] list: horizontal(row direction) layout, and is list's default layout.
- [x] list: vertical(column direction) layout.
- [x] list: simple in-cell list, element must be **scalar** type.
- [x] map: horizontal(row direction) layout.
- [x] map: vertical(column direction) layout, and is map's default layout.
- [x] map: unordered map or hash map.
- [ ] map: ordered map.
- [x] map: simple in-cell map, both key and value must be **scalar** type.
- [x] nesting: unlimited nesting of message, list, and map.

### Default Values
- [x] each scalar type's default value is same as protobuf
- [x] interger: 0 
- [x] float: 0.0 
- [x] bool: false 
- [x] string: ""
- [x] bytes: ""
- [x] in-cell message: each field's default value is same as protobuf 
- [x] in-cell list: element's default value is same as protobuf 
- [x] in-cell map: both key and value's default value is same as protobuf 
- [x] message: all fields of a message is empty(default values)
- [x] list: empty message will not be appended if list's element(message type) is empty
- [x] map: empty message will not be inserted if map's value(message type) is empty

### Merge
- [ ] merge multiple workbooks
- [ ] merge multiple worksheets

### Datetime
> [ISO 8601](https://www.wikiwand.com/en/ISO_8601) and [RFC 3339](https://tools.ietf.org/html/rfc3339)
> 
> [Understanding about RFC 3339 for Datetime and Timezone Formatting in Software Engineering](https://medium.com/easyread/understanding-about-rfc-3339-for-datetime-formatting-in-software-engineering-940aa5d5f68a)
> ```
> # This is acceptable in ISO 8601 and RFC 3339 (with T)
> 2019-10-12T07:20:50.52Z
> # This is only accepted in RFC 3339 (without T)
> 2019-10-12 07:20:50.52Z
> ```
> - "Z": stands for Zero timezone (UTC+0). Or equal to +00:00 in the RFC 3339.
> - RFC 3339 is following the ISO 8601 DateTime format. The only difference is RFC allows us to replace "T" with "space".

- [x] Timestamp: based on `google.protobuf.Timestamp`
- [ ] Timezone: time zones and time offsets
- [x] Datetime: format: `yyyy-MM-dd HH:mm:ss`, e.g.: `2020-01-01 05:10:00`
- [ ] Date: format: `yyyy-MM-dd`, e.g.: `2020-01-01`
- [ ] Time: format: `HH:mm:ss`, e.g.: `05:10:00`
- [x] Duration: based on`google.protobuf.Duration` 
- [x] Duration: format: `form "72h3m0.5s"`, see: [golang duration string form](https://golang.org/pkg/time/#Duration.String)
  
### Transpose
- [x] Interchange the rows and columns of a worksheet.

### Validation
- [ ] Min
- [ ] Max
- [ ] Range
- [ ] Options: e.g., enum type
- [ ] Foreign key

### Error Message
- [ ] report clear and precise error messages when converter failed, please refer to the programming language compiler
- [ ] use golang template to define error message template
- [ ] multiple languages support, focused on English and Simplified Chinese

### Performace
- [ ] stress test
- [ ] one goroutine process one row

## Types
- Scalar
- Message(struct)
- List
- Map(unordered)
- Timestamp
- Duration
