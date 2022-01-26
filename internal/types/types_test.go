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

func TestMatchBoringInteger(t *testing.T) {
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
			name: "integer 0",
			args: args{
				text: "0",
			},
			want: nil,
		},
		{
			name: "integer 1",
			args: args{
				text: "1",
			},
			want: nil,
		},
		{
			name: "integer 01",
			args: args{
				text: "01",
			},
			want: nil,
		},
		{
			name: "boring integer 0.0",
			args: args{
				text: "0.0",
			},
			want: []string{"0.0", "0"},
		},
		{
			name: "boring integer 1.000",
			args: args{
				text: "1.000",
			},
			want: []string{"1.000", "1"},
		},
		{
			name: "scientific-notation integer 1.0000001e7",
			args: args{
				text: "1.0000001e7",
			},
			want: nil,
		},
		{
			name: "scientific-notation integer 1.0000001E-7",
			args: args{
				text: "1.0000001E-7",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchBoringInteger(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchBoringInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchMap(t *testing.T) {
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
			name: "normal map",
			args: args{
				text: "map<int32, Server>",
			},
			want: []string{"map<int32, Server>", "int32", " Server"},
		},
		{
			name: "enum keyed map",
			args: args{
				text: "map<enum<.ServerType>, Server>",
			},
			want: []string{"map<enum<.ServerType>, Server>", "enum<.ServerType>", " Server"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MatchMap(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
