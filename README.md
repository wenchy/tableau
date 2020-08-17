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
- [ ] metatable: a message to describe the worksheet's metadata
- [ ] metafield: a message to describe the caption's metadata
- [x] captrow: caption row, exact row number of caption at worksheet
- [ ] descrow: exact row, number of description at wooksheet
- [x] datarow: data row, start row number of data

### Generator
- [ ] generate xlsx template by proto: **proto -> xlsx template**
- [ ] generate proto by xlsx template: **proto -> xlsx template**

### Conversion
- [x] xlsx -> JSON
- [x] xlsx -> protobin
- [x] xlsx -> prototext
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
- [x] message: multi-level nested message
- [x] message: simple in-cell message
- [x] list: horizontal(row direction) layout, and is list's default layout
- [x] list: vertical(column direction) layout
- [x] list: multi-level nested list
- [x] list: horizontal layout list
- [x] list: vertical layout list
- [x] list: simple in-cell list, element must be **scalar** type.
- [x] map: horizontal(row direction) layout
- [x] map: vertical(column direction) layout, and is map's default layout
- [x] map: multi-level nested map
- [x] map: unordered map or hash map
- [ ] map: ordered map
- [x] map: simple in-cell map, both key and value must be **scalar** type
- [ ] nested types: unlimited nesting of message, list, and map

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
- [ ] message: no empty message will be spawned if all fields of a message are empty

### Merge
- [ ] merge multiple workbooks
- [ ] merge multiple worksheets

### Time
- [x] Timestamp: `google.protobuf.Timestamp`
- [ ] Timestamp: timezone problem
- [x] Datetime: format: `yyyy-MM-dd HH:mm:ss`, based on Timestamp
- [ ] Date: format: `yyyy-MM-dd`, ignore day time based on Timestamp
- [ ] Time: format: `HH:mm:ss`
- [x] Duration: `google.protobuf.Duration` 
  
### Transpose
- [x] Interchange the rows and columns of a worksheet.

### Validation
- [ ] Min
- [ ] Max
- [ ] Range
- [ ] Options: e.g. enum type
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
