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
- [ ] Bidirectional conversion: **xlsx <-> json**.
- [ ] Merge of multi-level nested list, row direction
- [ ] Merge of multi-level nested list, column direction
- [x] Merge of multi-level nested map
- [ ] Ordered Map
- [ ] Merge of multiple workbooks or worksheets
- [ ] Timezone of type Timestamp
- [ ] Simple key-value configuration: flip worksheet 90° (degrees) to exchange row and column 
- [x] In cell array
- [] In cell map

## Types
- Scalar
- Timestamp
- Duration
- One-level List
- Multi-level List
- One-level Map(unordered)
- Multi-level Nested Map

