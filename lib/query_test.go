package lib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserExample struct {
	UserId     int     `json:"user_id"`
	Id         int     `json:"id"`
	FloatValue float64 `json:"float_value"`
	Path       string  `json:"-" path:"path"`
}

func TestBuildPath(t *testing.T) {
	type args struct {
		resourcePath string
		values       interface{}
	}
	tests := []struct {
		name string
		args args
		want string
		error
	}{
		{
			name: "_id",
			want: "users/3",
			args: args{
				"users/{user_id}",
				UserExample{UserId: 3},
			},
		},
		{
			name: "id",
			want: "users/3",
			args: args{
				"users/{id}",
				UserExample{Id: 3},
			},
		},
		{
			name: "id with large value",
			want: "users/3922799",
			args: args{
				"users/{id}",
				UserExample{Id: 3922799},
			},
		},
		{
			name: "no substitution",
			want: "users",
			args: args{
				"users",
				UserExample{},
			},
		},
		{
			name: "escaping",
			want: "root/a/%3F/c",
			args: args{
				"root/{path}",
				UserExample{Path: "a/?/c"}},
		},
		{
			name: "root path",
			want: "root/a/my-path",
			args: args{
				"root/{path}",
				UserExample{Path: "/a/my-path"}},
		},
		{
			name: "empty path",
			want: "root/",
			args: args{
				"root/{path}",
				UserExample{Path: ""}},
		},
		{
			name:  "validating int zero value",
			error: fmt.Errorf("missing required field: UserExample{}.id"),
			args: args{
				"root/{id}",
				UserExample{Id: 0}},
		},
		{
			name: "validating empty path value",
			want: "root/",
			args: args{
				"root/{path}",
				UserExample{}},
		},
		{
			name: "escaping spaces",
			want: "root/a%20file%20name.text",
			args: args{
				"root/{path}",
				UserExample{Path: "a file name.text"}},
		},
		{
			name: "given a map",
			want: "root/a/my-path",
			args: args{
				"root/{path}",
				map[string]interface{}{"path": "a/my-path"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildPath(tt.args.resourcePath, tt.args.values)
			if tt.error != nil {
				assert.Error(t, err, tt.Error())
			} else {
				assert.Equalf(t, tt.want, got, "BuildPath(%v, %v)", tt.args.resourcePath, tt.args.values)
			}
		})
	}
}
