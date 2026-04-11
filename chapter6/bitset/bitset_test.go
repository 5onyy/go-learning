package main

import (
	"testing"
)

func TestHas(t *testing.T) {
	tests := []struct {
		name  string
		add   []uint64
		check uint64
		want  bool
	}{
		{"present", []uint64{5, 64, 100}, 5, true},
		{"present boundary", []uint64{63}, 63, true},
		{"absent", []uint64{5, 64}, 10, false},
		{"empty set", []uint64{}, 0, false},
		{"exact bucket boundary", []uint64{64}, 64, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bs BitSet
			for _, v := range tt.add {
				bs.Add(v)
			}
			if got := bs.Has(tt.check); got != tt.want {
				t.Errorf("Has(%d) = %v, want %v", tt.check, got, tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		adds    []uint64
		wantLen int
		wantStr string
	}{
		{"single", []uint64{5}, 1, "{5}"},
		{"multiple", []uint64{2, 5, 64}, 3, "{2 5 64}"},
		{"duplicate ignored", []uint64{5, 5, 5}, 1, "{5}"},
		{"zero value", []uint64{0}, 1, "{0}"},
		{"large value", []uint64{127}, 1, "{127}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bs BitSet
			for _, v := range tt.adds {
				bs.Add(v)
			}
			if bs.Len() != tt.wantLen {
				t.Errorf("Len() = %d, want %d", bs.Len(), tt.wantLen)
			}
			if bs.String() != tt.wantStr {
				t.Errorf("String() = %q, want %q", bs.String(), tt.wantStr)
			}
		})
	}
}

func TestAddAll(t *testing.T) {
	tests := []struct {
		name    string
		vals    []int
		wantLen int
		wantStr string
	}{
		{"multiple values", []int{1, 2, 3}, 3, "{1 2 3}"},
		{"empty args", []int{}, 0, "{}"},
		{"with duplicates", []int{5, 5, 10}, 2, "{5 10}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bs BitSet
			bs.AddAll(tt.vals...)
			if bs.Len() != tt.wantLen {
				t.Errorf("Len() = %d, want %d", bs.Len(), tt.wantLen)
			}
			if bs.String() != tt.wantStr {
				t.Errorf("String() = %q, want %q", bs.String(), tt.wantStr)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name    string
		add     []uint64
		remove  uint64
		wantLen int
		wantStr string
	}{
		{"remove existing", []uint64{5, 10}, 5, 1, "{10}"},
		{"remove absent (no-op)", []uint64{5, 10}, 99, 2, "{5 10}"},
		{"remove only element", []uint64{42}, 42, 0, "{}"},
		{"remove across bucket boundary", []uint64{63, 64}, 64, 1, "{63}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bs BitSet
			for _, v := range tt.add {
				bs.Add(v)
			}
			bs.Remove(tt.remove)
			if bs.Len() != tt.wantLen {
				t.Errorf("Len() = %d, want %d", bs.Len(), tt.wantLen)
			}
			if bs.String() != tt.wantStr {
				t.Errorf("String() = %q, want %q", bs.String(), tt.wantStr)
			}
		})
	}
}

func TestClear(t *testing.T) {
	tests := []struct {
		name string
		add  []uint64
	}{
		{"non-empty set", []uint64{1, 2, 3, 100}},
		{"already empty", []uint64{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bs BitSet
			for _, v := range tt.add {
				bs.Add(v)
			}
			bs.Clear()
			if bs.Len() != 0 {
				t.Errorf("Len() = %d, want 0 after Clear", bs.Len())
			}
			if bs.String() != "{}" {
				t.Errorf("String() = %q, want \"{}\" after Clear", bs.String())
			}
		})
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		name    string
		x       []uint64
		y       []uint64
		wantLen int
		wantStr string
	}{
		{"disjoint sets", []uint64{1, 2}, []uint64{3, 4}, 4, "{1 2 3 4}"},
		{"overlapping sets", []uint64{1, 2, 3}, []uint64{2, 3, 4}, 4, "{1 2 3 4}"},
		{"union with empty", []uint64{1, 2}, []uint64{}, 2, "{1 2}"},
		{"both empty", []uint64{}, []uint64{}, 0, "{}"},
		{"nil other (no-op)", []uint64{1, 2}, nil, 2, "{1 2}"},
		{"across bucket boundary", []uint64{63}, []uint64{64}, 2, "{63 64}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x BitSet
			for _, v := range tt.x {
				x.Add(v)
			}
			var other *BitSet
			if tt.y != nil {
				var y BitSet
				for _, v := range tt.y {
					y.Add(v)
				}
				other = &y
			}
			x.Union(other)
			if x.Len() != tt.wantLen {
				t.Errorf("Len() = %d, want %d", x.Len(), tt.wantLen)
			}
			if x.String() != tt.wantStr {
				t.Errorf("String() = %q, want %q", x.String(), tt.wantStr)
			}
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		name    string
		x       []uint64
		y       []uint64
		wantLen int
		wantStr string
	}{
		{"common elements", []uint64{1, 2, 3}, []uint64{2, 3, 4}, 2, "{2 3}"},
		{"no common elements", []uint64{1, 2}, []uint64{3, 4}, 0, "{}"},
		{"intersect with empty", []uint64{1, 2}, []uint64{}, 0, "{}"},
		{"nil other (no-op)", []uint64{1, 2}, nil, 2, "{1 2}"},
		{"identical sets", []uint64{5, 10}, []uint64{5, 10}, 2, "{5 10}"},
		{"across bucket boundary", []uint64{63, 64, 65}, []uint64{63, 65}, 2, "{63 65}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x BitSet
			for _, v := range tt.x {
				x.Add(v)
			}
			var other *BitSet
			if tt.y != nil {
				var y BitSet
				for _, v := range tt.y {
					y.Add(v)
				}
				other = &y
			}
			x.Intersect(other)
			if x.Len() != tt.wantLen {
				t.Errorf("Len() = %d, want %d", x.Len(), tt.wantLen)
			}
			if x.String() != tt.wantStr {
				t.Errorf("String() = %q, want %q", x.String(), tt.wantStr)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	t.Run("copy is independent", func(t *testing.T) {
		var x BitSet
		x.Add(1)
		x.Add(2)
		clone := x.Copy()
		clone.Add(99)
		if x.Has(99) {
			t.Error("mutating clone affected original")
		}
		if !clone.Has(99) {
			t.Error("clone missing added element")
		}
	})
	t.Run("copy preserves contents", func(t *testing.T) {
		var x BitSet
		x.Add(5)
		x.Add(100)
		clone := x.Copy()
		if clone.String() != x.String() {
			t.Errorf("clone.String() = %q, want %q", clone.String(), x.String())
		}
		if clone.Len() != x.Len() {
			t.Errorf("clone.Len() = %d, want %d", clone.Len(), x.Len())
		}
	})
	t.Run("nil receiver returns nil", func(t *testing.T) {
		var bs *BitSet
		if bs.Copy() != nil {
			t.Error("Copy() on nil receiver should return nil")
		}
	})
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name    string
		x       []uint64
		y       []uint64
		wantLen int
		wantStr string
	}{
		{"removes common elements", []uint64{1, 2, 3, 4}, []uint64{2, 3}, 2, "{1 4}"},
		{"no overlap, x unchanged", []uint64{1, 2}, []uint64{3, 4}, 2, "{1 2}"},
		{"y is superset, x becomes empty", []uint64{2, 3}, []uint64{1, 2, 3, 4}, 0, "{}"},
		{"y larger bucket range, x tail preserved", []uint64{1, 64}, []uint64{1}, 1, "{64}"},
		{"nil other (no-op)", []uint64{1, 2}, nil, 2, "{1 2}"},
		{"both empty", []uint64{}, []uint64{}, 0, "{}"},
		{"x empty", []uint64{}, []uint64{1, 2}, 0, "{}"},
		{"across bucket boundary", []uint64{63, 64}, []uint64{64}, 1, "{63}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x BitSet
			for _, v := range tt.x {
				x.Add(v)
			}
			var other *BitSet
			if tt.y != nil {
				var y BitSet
				for _, v := range tt.y {
					y.Add(v)
				}
				other = &y
			}
			x.Difference(other)
			if x.Len() != tt.wantLen {
				t.Errorf("Len() = %d, want %d", x.Len(), tt.wantLen)
			}
			if x.String() != tt.wantStr {
				t.Errorf("String() = %q, want %q", x.String(), tt.wantStr)
			}
		})
	}
}

func TestSymetricDifference(t *testing.T) {
	tests := []struct {
		name    string
		x       []uint64
		y       []uint64
		wantLen int
		wantStr string
	}{
		{"disjoint sets (union)", []uint64{1, 2}, []uint64{3, 4}, 4, "{1 2 3 4}"},
		{"overlapping sets", []uint64{1, 2, 3}, []uint64{2, 3, 4}, 2, "{1 4}"},
		{"identical sets (empty result)", []uint64{1, 2, 3}, []uint64{1, 2, 3}, 0, "{}"},
		{"nil other (no-op)", []uint64{1, 2}, nil, 2, "{1 2}"},
		{"x empty", []uint64{}, []uint64{1, 2}, 2, "{1 2}"},
		{"both empty", []uint64{}, []uint64{}, 0, "{}"},
		{"across bucket boundary", []uint64{63, 64}, []uint64{64, 65}, 2, "{63 65}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x BitSet
			for _, v := range tt.x {
				x.Add(v)
			}
			var other *BitSet
			if tt.y != nil {
				var y BitSet
				for _, v := range tt.y {
					y.Add(v)
				}
				other = &y
			}
			x.SymetricDifference(other)
			if x.Len() != tt.wantLen {
				t.Errorf("Len() = %d, want %d", x.Len(), tt.wantLen)
			}
			if x.String() != tt.wantStr {
				t.Errorf("String() = %q, want %q", x.String(), tt.wantStr)
			}
		})
	}
}

func TestElems(t *testing.T) {
	tests := []struct {
		name string
		add  []uint64
		want []uint64
	}{
		{"empty set", []uint64{}, []uint64{}},
		{"single element", []uint64{5}, []uint64{5}},
		{"multiple elements sorted", []uint64{100, 5, 64}, []uint64{5, 64, 100}},
		{"zero element", []uint64{0}, []uint64{0}},
		{"across bucket boundary", []uint64{63, 64}, []uint64{63, 64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bs BitSet
			for _, v := range tt.add {
				bs.Add(v)
			}
			got := bs.Elems()
			if len(got) != len(tt.want) {
				t.Fatalf("Elems() len = %d, want %d: got %v", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Elems()[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}

	t.Run("nil receiver returns nil", func(t *testing.T) {
		var bs *BitSet
		if bs.Elems() != nil {
			t.Error("Elems() on nil receiver should return nil")
		}
	})

	t.Run("matches Len()", func(t *testing.T) {
		var bs BitSet
		bs.AddAll(1, 2, 63, 64, 100)
		if len(bs.Elems()) != bs.Len() {
			t.Errorf("len(Elems()) = %d, Len() = %d; must match", len(bs.Elems()), bs.Len())
		}
	})
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		add  []uint64
		want string
	}{
		{"empty", []uint64{}, "{}"},
		{"single element", []uint64{7}, "{7}"},
		{"multiple elements sorted", []uint64{100, 5, 64}, "{5 64 100}"},
		{"zero element", []uint64{0}, "{0}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bs BitSet
			for _, v := range tt.add {
				bs.Add(v)
			}
			if got := bs.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
