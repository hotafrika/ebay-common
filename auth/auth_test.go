package auth

import "testing"

func Test_createB64String(t *testing.T) {
	type args struct {
		clientID     string
		clientSecret string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{
				clientID:     "user",
				clientSecret: "password",
			},
			want: "Basic dXNlcjpwYXNzd29yZA==",
		},
		{
			name: "invalid",
			args: args{
				clientID:     "123",
				clientSecret: "123",
			},
			want: "Basic MTIzOjEyMw==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createB64String(tt.args.clientID, tt.args.clientSecret); got != tt.want {
				t.Errorf("createB64String() = %v, want %v", got, tt.want)
			}
		})
	}
}
