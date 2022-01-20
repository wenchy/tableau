package types

import (
	"reflect"
	"testing"
)

func TestMatchList(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "scalar list",
			args: args{
				text: "[]uint32",
			},
			want: []string{"[]uint32", "", "uint32"},
		},
		{
			name: "struct list",
			args: args{
				text: "[Type]uint32",
			},
			want: []string{"[Type]uint32", "Type", "uint32"},
		},
		{
			name: "keyed struct list",
			args: args{
				text: "[Type]<uint32>",
			},
			want: []string{"[Type]<uint32>", "Type", "<uint32>"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchList(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchKeyedList(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "keyed struct list",
			args: args{
				text: "[Type]<uint32>",
			},
			want: []string{"[Type]<uint32>", "Type", "uint32"},
		},
		{
			name: "normal struct list",
			args: args{
				text: "[Type]uint32",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchKeyedList(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchKeyedList() = %v, want %v", got, tt.want)
			}
		})
	}
}
