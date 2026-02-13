package service

import "testing"

func TestParseAheadBehind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		wantBehind int
		wantAhead  int
		wantErr    bool
	}{
		{name: "valid", input: "2 3", wantBehind: 2, wantAhead: 3},
		{name: "valid with spaces", input: "  0   1 ", wantBehind: 0, wantAhead: 1},
		{name: "invalid token count", input: "2", wantErr: true},
		{name: "invalid behind", input: "x 1", wantErr: true},
		{name: "invalid ahead", input: "1 x", wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotBehind, gotAhead, err := ParseAheadBehind(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotBehind != tc.wantBehind || gotAhead != tc.wantAhead {
				t.Fatalf("got behind/ahead %d/%d, want %d/%d", gotBehind, gotAhead, tc.wantBehind, tc.wantAhead)
			}
		})
	}
}
