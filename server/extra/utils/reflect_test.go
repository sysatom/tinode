package utils

import "testing"

func foo1() {}
func foo2(a, b string) string {
	return ""
}
func foo3(a, b string, c int) *testing.InternalExample {
	return nil
}

func TestGetFunctionName(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case1",
			args: args{
				i: foo1,
			},
			want: "foo1",
		},
		{
			name: "case2",
			args: args{
				i: foo2,
			},
			want: "foo2",
		},
		{
			name: "case3",
			args: args{
				i: foo3,
			},
			want: "foo3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFunctionName(tt.args.i); got != tt.want {
				t.Errorf("GetFunctionName() = %v, want %v", got, tt.want)
			}
		})
	}
}
