package utils

import (
	"testing"
)

func TestJWT(t *testing.T) {
	type args struct {
		username  string
		uid       int
		secretKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				username:  "zineyu",
				uid:       1,
				secretKey: "zineyu",
			},
			wantErr: false,
		},
		{
			name: "Test 2",
			args: args{
				username:  "zine",
				uid:       2,
				secretKey: "zineyu",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateJWT(tt.args.username, tt.args.uid, tt.args.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			parse, err := ValidJWT(got, tt.args.secretKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidJWT() error = %v", err)
				return
			}

			if parse.Username != tt.args.username || parse.Uid != tt.args.uid {
				t.Errorf("want: %+v\ngot: %+v", tt.args, parse)
				return
			}
		})
	}
}
