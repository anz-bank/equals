package equals

import "testing"

func TestListHas(t *testing.T) {
	type Foo struct {
		A []string
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
			element: Foo{A: []string{"a", "b", "c"}},
			list: []Foo{{
				A: []string{"a", "b", "c"},
			}},
			want: true,
		},
		{
			element: Foo{A: []string{"b", "c"}},
			list: []Foo{{
				A: []string{"a", "b", "c"},
			}},
			want: false,
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
