# Metadata Conversion of Workbook to Protoconf
convert workbook to protoconf.

## Rules: basic

workbook: `DemoTest(test)`, worksheet: `DemoActivity(activity)`

- protoconf file name is `test.proto`. If with no `()`, name will be `demo_test.proto`
- configuration message name is `Activity`. If with no `()`, name will be `DemoActivity`
- list: `TYPE[]ELEM`,  TYPE is this column type, ELEM is message name(must not conflict with the protobuf keyword).
- map: `map[KEY]VALUE`, KEY must be scalar types, and VALUE is message name(must not conflict with build-in scalar type).
- common message types: `(TYPE)Name`, TYPE muse be message defined in `common.proto`. In convention, `common.Item` is equal to `Item` if not specify the prefix `common` which is a proto file name.

| ActivityID          | ActivityName | ChapterID          | ChapterName | SectionID       | SectionName | (common.Item)SectionItem1ID | (common.Item)SectionItem1Num | (common.Item)SectionItem2ID | (common.Item)SectionItem2Num |
| ------------------- | ------------ | ------------------ | ----------- | --------------- | ----------- | --------------------------- | ---------------------------- | --------------------------- | ---------------------------- |
| map[uint32]Activity | string       | map[uint32]Chapter | string      | uint32[]Section | int32       | int32                       | int32                        | int32                       | int32                        |
| 1                   | activity1    | 1                  | chapter1    | 1               | section1    | 1001                        | 1                            | 1002                        | 2                            |
| 1                   | activity1    | 1                  | chapter1    | 2               | section2    | 1001                        | 1                            | 1002                        | 2                            |
| 1                   | activity1    | 2                  | chapter2    | 1               | section1    | 1001                        | 1                            | 1002                        | 2                            |
| 2                   | activity2    | 1                  | chapter1    | 1               | section1    | 1001                        | 1                            | 100)2                       | 2                            |

```
// common.proto
message Item {
	int32 id = 1 [(caption) = "ID"];
	int32 num= 1 [(caption) = "Num"];
}
```

without prefix:
```
// test.proto
message Activity{
	map<uint32, Activity> activity_map = 1 [(key) = "ActivityID"];
	message Activity {
		uint32 id= 1 [(caption) = "ActivityID"];
		string name = 2 [(caption) = "ActivityName "];
		map<uint32, Chapter> chapter_map = 1 [(key) = "ChapterID"];
	}
	message Chapter {
		uint32 id= 1 [(caption) = "ChapterID"];
		string name = 2 [(caption) = "ChapterName"];
		repeated Section section_list = 3 [(layout) = COMPOSITE_LAYOUT_VERTICAL];
	}
	message Section {
		uint32 id= 1 [(caption) = "SectionID"];
		string name = 1 [(caption) = "SectionName"];
		repeated Item item_list = 3 [(caption) = "SectionItem"];
	}
}
```

with prefix: 
```
// test.proto
message Activity {
	map<uint32, Activity> activity_map = 1 [(key) = "ActivityID"];
	message Activity {
		uint32 activity_id= 1 [(caption) = "ActivityID"];
		string activity_name = 2 [(caption) = "ActivityName "];
		map<uint32, Chapter> chapter_map = 1 [(key) = "ChapterID"];
	}
	message Chapter {
		uint32 chapter_id= 1 [(caption) = "ChapterID"];
		string chapter_name = 2 [(caption) = "ChapterName"];
		repeated Section section_list = 3 [(layout) = COMPOSITE_LAYOUT_VERTICAL];
	}
	message Section {
		uint32 section_id= 1 [(caption) = "SectionID"];
		string section_name = 1 [(caption) = "SectionName"];
		repeated Item section_item_list = 3 [(caption) = "SectionItem"];
	}
}
```

## Rules: in-cell

workbook: `DemoTest(test)`, worksheet: `Environment(env)`

| ID     | Name   | InCellMessage                       | InCellList | InCellMap        | InCellMessageList            | InCellMessageMap                      |
| ------ | ------ | ----------------------------------- | ---------- | ---------------- | ---------------------------- | ------------------------------------- |
| uint32 | string | {int32 id,string desc,uint32 value} | []int32    | map[int32]string | []Elem{int32 id,string desc} | map[int32]Value{int32 id,string desc} |
| 1      | Earth  | 1,desc,100                          | 1,2,3      | 1:hello,2:world  | {1,hello},{2,world}          | 1:{1,hello},2:{2,world}               |

```
// test.proto
message Env {
	uint32 ID = 1 [(caption) = "ID"];
	string name = 2 [(caption) = "Name "];
	InCellMessage in_cell_message= 3 [(caption) = "InCellMessage"];
	repeated int32 in_cell_list= 4 [(caption) = "InCellList"];
	map<int32, string> in_cell_map = 5 [(caption) = "InCellMap"];
	repeated Elem in_cell_message_list= 5 [(caption) = "InCellMessageList"];
    map<int32, Value> in_cell_message_map = 6 [(caption) = "InCellMessageMap"];

	message InCellMessage {
		int32 id = 1;
		string desc= 2; // defaut name: field + [tagid]
		uint32 value= 3;
	}
    message Elem {
		int32 id = 1;
		string desc= 2; // defaut name: field + [tagid]
	}
    message Value {
		int32 id = 1;
		string desc= 2; // defaut name: field + [tagid]
	}
}
```

- in-cell message: comma separeted sequence: `{TYPE [NAME],TYPE [NAME]}`, NAME is optional, and will be auto generated as `field + [tagid]` if not specified.
- in-cell list: `[]TYPE`, TYPE must be scalar type.
- in-cell map: `map[KEY]VALUE`, KEY and VALUE must be scalar types.
- in-cell message list: `[]TYPE`, TYPE must be message type.
- in-cell message map: `map[KEY]VALUE`, KEY is scalar, and VALUE must be message type.