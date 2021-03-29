package equals

import "testing"

func TestListHas(t *testing.T) {
	type Foo struct {
		a []string
	}
	tests := []struct {
		name    string
		element interface{}
		list    interface{}
		want    bool
	}{
		{
			element: "a",
			list:    []string{"a", "b", "c"},
			want:    true,
		},
		{
			element: Foo{a: []string{"a", "b", "c"}},
			list:    []Foo{{
				a: []string{"a", "b", "c"},
			}},
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListHas(tt.element, tt.list); got != tt.want {
				t.Errorf("ElementInListRec() = %v, want %v", got, tt.want)
			}
		})
	}
}
