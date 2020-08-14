# Tableau
A powerful configuration convertion tool based on protobuf.

## Design
- Convert **xlsx** to **json**, json is the first-class citizen of exporting targets.
- Comments in json: add one more comment key-value pair, and the comment key is prefixed with "#".
- Use **protobuf** as the IDL(Interface Description Language) to define the structure of **xlsx**.
- Use **golang** to develop the conversion engine.
- Multiple languages support, thanks to **protobuf**.

## Concept
- Importer: xlsx importer
- IR: Intermediate Representation, use proto-bin.
- Filter: filter the IR.
- Exporter: json exporter, proto-bin exporter, proto-text exporter, xml exporter, sqlite3 exporter, and so on.

## TODO
- [ ] Battle-tested of different languages: Golang, C#, Javascript/Typescript, C++ and so on.
- [ ] Generate xlsx template by proto: **proto -> xlsx template**.
- [ ] Generate proto by xlsx template: **proto -> xlsx template**.
- [x] Convert xlsx to json: **xlsx <-> json**.
- [ ] Convert json to xlsx: **json <-> xlsx**.
- [x] List: multi-level nested list
- [x] List: horizontal layout list
- [x] List: vertical layout list
- [x] List: simple in-cell list, element must be **scalar**.
- [x] Map: multi-level nested map
- [x] Map: unordered map or hash map
- [ ] Map: ordered map
- [x] Map: simple in-cell map, both key and value must be **scalar**.
- [ ] Merge: multiple workbooks merge
- [ ] Merge: multiple worksheets merge
- [x] Timestamp: `google.protobuf.Timestamp`
- [x] Timestamp: timezone problem
- [x] Datetime: format: `yyyy-MM-dd HH:mm:ss`, based on Timestamp
- [ ] Date: format: `yyyy-MM-dd`, ignore day time based on Timestamp
- [ ] Time: format: `HH:mm:ss`
- [ ] Simple key-value configuration: flip worksheet 90° (degrees) to exchange row and column 
## Types
- Scalar
- Timestamp
- Duration
- One-level List
- Multi-level List
- One-level Map(unordered)
- Multi-level Nested Map

