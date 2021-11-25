package protogen

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type sheetHeader struct {
	namerow []string
	typerow []string
	noterow []string
}

func getCell(row []string, cursor int) string {
	return strings.TrimSpace(row[cursor])
}
func (sh *sheetHeader) getNameCell(cursor int) string {
	return getCell(sh.namerow, cursor)
}

func (sh *sheetHeader) getTypeCell(cursor int) string {
	return getCell(sh.typerow, cursor)
}
func (sh *sheetHeader) getNoteCell(cursor int) string {
	return getCell(sh.noterow, cursor)
}

type sheet struct {
	ws             *tableaupb.Worksheet
	writer         *bufio.Writer
	isLastSheet    bool
	nestedMessages map[string]*tableaupb.Field // type name -> field
}

func (s *sheet) export() error {
	s.writer.WriteString(fmt.Sprintf("message %s {\n", s.ws.Name))
	s.writer.WriteString(fmt.Sprintf("  option (tableau.worksheet) = {%s};\n", genPrototext(s.ws.Options)))
	s.writer.WriteString("\n")
	// generate the fields
	depth := 1
	for i, field := range s.ws.Fields {
		tagid := i + 1
		if err := s.exportField(depth, tagid, field); err != nil {
			return err
		}
	}
	s.writer.WriteString("}\n")
	if !s.isLastSheet {
		s.writer.WriteString("\n")
	}
	return nil
}

func (s *sheet) exportField(depth int, tagid int, field *tableaupb.Field) error {
	head := "%s%s"
	if field.Card != "" {
		head += " " // cardinality exists
	}
	s.writer.WriteString(fmt.Sprintf(head+"%s %s = %d [(tableau.field) = {%s}];\n", indent(depth), field.Card, field.Type, field.Name, tagid, genPrototext(field.Options)))

	if !field.TypeDefined && field.Fields != nil {
		// iff field is a map or list and message type is not imported.
		nestedMsgName := field.Type
		if field.MapEntry != nil {
			nestedMsgName = field.MapEntry.ValueType
		}

		if isSameFieldMessageType(field, s.nestedMessages[nestedMsgName]) {
			// if the nested message is the same as the previous one,
			// just use the previous one, and don't generate a new one.
			return nil
		}

		// bookkeeping this nested msessage, so we can check if we can reuse it later.
		s.nestedMessages[nestedMsgName] = field

		s.writer.WriteString("\n")
		s.writer.WriteString(fmt.Sprintf("%smessage %s {\n", indent(depth), nestedMsgName))
		for i, f := range field.Fields {
			tagid := i + 1
			if err := s.exportField(depth+1, tagid, f); err != nil {
				return err
			}
		}
		s.writer.WriteString(fmt.Sprintf("%s}\n", indent(depth)))
	}
	return nil
}

func genPrototext(m protoreflect.ProtoMessage) string {
	// text := proto.CompactTextString(field.Options)
	bin, err := prototext.Marshal(m)
	if err != nil {
		panic(err)
	}
	// NOTE: remove redundant spaces/whitespace from a string
	// refer: https://stackoverflow.com/questions/37290693/how-to-remove-redundant-spaces-whitespace-from-a-string-in-golang
	text := strings.Join(strings.Fields(string(bin)), " ")
	return text
}

func indent(depth int) string {
	return strings.Repeat("  ", depth)
}

func isSameFieldMessageType(left, right *tableaupb.Field) bool {
	if left == nil || right == nil {
		return false
	}
	if left.Fields == nil || right.Fields == nil {
		return false
	}
	if len(left.Fields) != len(right.Fields) ||
		left.Type != right.Type ||
		left.Card != right.Card {
		return false
	}

	for i, l := range left.Fields {
		r := right.Fields[i]
		if !proto.Equal(l, r) {
			return false
		}
	}
	return true
}
