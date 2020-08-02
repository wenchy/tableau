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

