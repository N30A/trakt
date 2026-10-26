package validate

import "testing"

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "valid email",
			input:   "john@example.com",
			want:    "john@example.com",
			wantErr: false,
		},
		{
			name:    "trims whitespace",
			input:   "  john@example.com  ",
			want:    "john@example.com",
			wantErr: false,
		},
		{
			name:    "rejects display name",
			input:   "John Doe <john@example.com>",
			wantErr: true,
		},
		{
			name:    "valid email with first and last name",
			input:   "john.doe@example.com",
			want:    "john.doe@example.com",
			wantErr: false,
		},
		{
			name:    "invalid email",
			input:   "not-an-email",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Email(test.input)

			if err != nil && !test.wantErr {
				t.Fatalf("Email() error = %v, wantErr %v", err, test.wantErr)
			}

			if err == nil && test.wantErr {
				t.Fatalf("Email() error = nil, wantErr %v", test.wantErr)
			}

			if got != test.want {
				t.Errorf("Email() = %q, want %q", got, test.want)
			}
		})
	}
}
