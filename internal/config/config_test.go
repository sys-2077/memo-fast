package config

import "testing"

func TestNormalizeCollectionName_IdempotentPrefix(t *testing.T) {
	got := NormalizeCollectionName("memo_foo")
	if got != "memo_foo" {
		t.Fatalf("expected memo_foo, got %q", got)
	}
}

func TestNormalizeCollectionName_FormatsProjectNames(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "space", in: "my project", want: "memo_myproject"},
		{name: "dash", in: "my-project", want: "memo_my_project"},
		{name: "dot", in: "my.project", want: "memo_my_project"},
		{name: "mixed case prefixed", in: "Memo_Foo", want: "memo_foo"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeCollectionName(tc.in)
			if got != tc.want {
				t.Fatalf("NormalizeCollectionName(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
