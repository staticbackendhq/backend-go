package backend

import "testing"

func TestUsersPath(t *testing.T) {
	tests := []struct {
		name  string
		roles []int
		want  string
	}{
		{
			name: "without role",
			want: "/account/users",
		},
		{
			name:  "with role",
			roles: []int{25},
			want:  "/account/users?role=25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := usersPath(tt.roles...); got != tt.want {
				t.Fatalf("expected %s got %s", tt.want, got)
			}
		})
	}
}
