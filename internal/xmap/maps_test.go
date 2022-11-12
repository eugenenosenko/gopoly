package xmap

import (
	"reflect"
	"testing"
)

func TestDifference(t *testing.T) {
	type args[M interface{ map[string]string }] struct {
		m1 M
		m2 M
	}
	type testCase[M interface{ map[string]string }] struct {
		name string
		args args[M]
		want M
	}
	tests := []testCase[map[string]string]{
		{
			name: "should correctly generate difference for different maps with second map having priority #1",
			args: args[map[string]string]{
				m1: map[string]string{"a": "1", "b": "2", "c": "33"},
				m2: map[string]string{"a": "11", "b": "22", "c": "33", "d": "44"},
			},
			want: map[string]string{"a": "11", "b": "22", "d": "44"},
		},
		{
			name: "should correctly generate difference for different maps with second map having priority #2",
			args: args[map[string]string]{
				m1: map[string]string{"a": "11", "b": "22", "c": "33", "d": "44"},
				m2: map[string]string{"a": "1", "b": "2", "c": "33"},
			},
			want: map[string]string{"a": "1", "b": "2", "d": "44"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Difference(tt.args.m1, tt.args.m2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}
