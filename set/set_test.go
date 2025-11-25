package set

import "testing"

func TestHashSetBasicOperations(t *testing.T) {
	set := NewHashSet[int]()

	if !set.Add(1) {
		t.Error("Expected Add to return true for new element")
	}
	if set.Add(1) {
		t.Error("Expected Add to return false for duplicate element")
	}

	if !set.Contains(1) {
		t.Error("Expected set to contain 1")
	}
	if set.Contains(2) {
		t.Error("Expected set not to contain 2")
	}

	set.Add(2)
	set.Add(3)
	if set.Size() != 3 {
		t.Errorf("Expected size 3, got %d", set.Size())
	}

	if !set.Remove(2) {
		t.Error("Expected Remove to return true")
	}
	if set.Remove(2) {
		t.Error("Expected Remove to return false for non-existent element")
	}
	if set.Size() != 2 {
		t.Errorf("Expected size 2 after remove, got %d", set.Size())
	}

	if set.IsEmpty() {
		t.Error("Expected set not to be empty")
	}

	set.Clear()
	if !set.IsEmpty() {
		t.Error("Expected set to be empty after Clear")
	}
	if set.Size() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", set.Size())
	}
}

func TestHashSetSetOperations(t *testing.T) {
	set1 := NewHashSet[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2 := NewHashSet[int]()
	set2.Add(3)
	set2.Add(4)
	set2.Add(5)

	union := set1.Union(set2)
	if union.Size() != 5 {
		t.Errorf("Expected union size 5, got %d", union.Size())
	}
	for i := 1; i <= 5; i++ {
		if !union.Contains(i) {
			t.Errorf("Expected union to contain %d", i)
		}
	}

	intersection := set1.Intersection(set2)
	if intersection.Size() != 1 {
		t.Errorf("Expected intersection size 1, got %d", intersection.Size())
	}
	if !intersection.Contains(3) {
		t.Error("Expected intersection to contain 3")
	}

	diff := set1.Difference(set2)
	if diff.Size() != 2 {
		t.Errorf("Expected difference size 2, got %d", diff.Size())
	}
	if !diff.Contains(1) || !diff.Contains(2) {
		t.Error("Expected difference to contain 1 and 2")
	}
	if diff.Contains(3) {
		t.Error("Expected difference not to contain 3")
	}
}

func TestTreeSetBasicOperations(t *testing.T) {
	set := NewTreeSet[int](func(a, b int) bool { return a < b })

	set.Add(5)
	set.Add(2)
	set.Add(8)
	set.Add(1)
	set.Add(9)

	slice := set.ToSlice()
	expected := []int{1, 2, 5, 8, 9}
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("Expected ordered element %d at position %d, got %d", expected[i], i, v)
		}
	}

	if set.Add(5) {
		t.Error("Expected Add to return false for duplicate")
	}

	if !set.Contains(8) {
		t.Error("Expected set to contain 8")
	}
	if set.Contains(7) {
		t.Error("Expected set not to contain 7")
	}

	set.Remove(2)
	slice = set.ToSlice()
	expected = []int{1, 5, 8, 9}
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("After remove, expected %d at position %d, got %d", expected[i], i, v)
		}
	}

	if set.Size() != 4 {
		t.Errorf("Expected size 4, got %d", set.Size())
	}
}

func TestTreeSetWithStrings(t *testing.T) {
	set := NewTreeSet[string](func(a, b string) bool { return a < b })

	set.Add("banana")
	set.Add("apple")
	set.Add("cherry")
	set.Add("date")

	slice := set.ToSlice()
	expected := []string{"apple", "banana", "cherry", "date"}
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("Expected %s at position %d, got %s", expected[i], i, v)
		}
	}
}

func TestTreeSetSetOperations(t *testing.T) {
	less := func(a, b int) bool { return a < b }
	set1 := NewTreeSet[int](less)
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2 := NewTreeSet[int](less)
	set2.Add(3)
	set2.Add(4)
	set2.Add(5)

	union := set1.Union(set2)
	slice := union.ToSlice()
	expected := []int{1, 2, 3, 4, 5}
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("Expected union element %d at position %d, got %d", expected[i], i, v)
		}
	}

	intersection := set1.Intersection(set2)
	if intersection.Size() != 1 || !intersection.Contains(3) {
		t.Error("Expected intersection to contain only 3")
	}
}

func TestLinkedHashSetBasicOperations(t *testing.T) {
	set := NewLinkedHashSet[int]()

	set.Add(5)
	set.Add(2)
	set.Add(8)
	set.Add(1)

	slice := set.ToSlice()
	expected := []int{5, 2, 8, 1}
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("Expected element %d at position %d, got %d", expected[i], i, v)
		}
	}

	if !set.Contains(8) {
		t.Error("Expected set to contain 8")
	}

	set.Remove(2)
	slice = set.ToSlice()
	expected = []int{5, 8, 1}
	for i, v := range slice {
		if v != expected[i] {
			t.Errorf("After remove, expected %d at position %d, got %d", expected[i], i, v)
		}
	}

	if set.Size() != 3 {
		t.Errorf("Expected size 3, got %d", set.Size())
	}

	if set.Add(5) {
		t.Error("Expected Add to return false for duplicate")
	}
}

func TestLinkedHashSetRemoveHeadAndTail(t *testing.T) {
	set := NewLinkedHashSet[string]()
	set.Add("first")
	set.Add("middle")
	set.Add("last")

	set.Remove("first")
	slice := set.ToSlice()
	if len(slice) != 2 || slice[0] != "middle" || slice[1] != "last" {
		t.Error("Failed to remove head correctly")
	}

	set.Remove("last")
	slice = set.ToSlice()
	if len(slice) != 1 || slice[0] != "middle" {
		t.Error("Failed to remove tail correctly")
	}

	set.Remove("middle")
	if !set.IsEmpty() {
		t.Error("Expected set to be empty")
	}
}

func TestLinkedHashSetSetOperations(t *testing.T) {
	set1 := NewLinkedHashSet[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2 := NewLinkedHashSet[int]()
	set2.Add(3)
	set2.Add(4)
	set2.Add(5)

	union := set1.Union(set2)
	if union.Size() != 5 {
		t.Errorf("Expected union size 5, got %d", union.Size())
	}

	intersection := set1.Intersection(set2)
	if intersection.Size() != 1 || !intersection.Contains(3) {
		t.Error("Expected intersection to contain only 3")
	}

	diff := set1.Difference(set2)
	slice := diff.ToSlice()
	if len(slice) != 2 || slice[0] != 1 || slice[1] != 2 {
		t.Error("Expected difference to contain 1 and 2 in order")
	}
}

func BenchmarkHashSetAdd(b *testing.B) {
	set := NewHashSet[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Add(i)
	}
}

func BenchmarkTreeSetAdd(b *testing.B) {
	set := NewTreeSet[int](func(a, b int) bool { return a < b })
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Add(i)
	}
}

func BenchmarkLinkedHashSetAdd(b *testing.B) {
	set := NewLinkedHashSet[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Add(i)
	}
}

func BenchmarkHashSetContains(b *testing.B) {
	set := NewHashSet[int]()
	for i := 0; i < 1000; i++ {
		set.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Contains(i % 1000)
	}
}

func BenchmarkTreeSetContains(b *testing.B) {
	set := NewTreeSet[int](func(a, b int) bool { return a < b })
	for i := 0; i < 1000; i++ {
		set.Add(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set.Contains(i % 1000)
	}
}

// Additional edge case tests
func TestEmptySets(t *testing.T) {
	hashSet := NewHashSet[int]()
	treeSet := NewTreeSet[int](func(a, b int) bool { return a < b })
	linkedSet := NewLinkedHashSet[int]()

	sets := []Settable[int]{hashSet, treeSet, linkedSet}
	names := []string{"HashSet", "TreeSet", "LinkedHashSet"}

	for i, set := range sets {
		if !set.IsEmpty() {
			t.Errorf("%s: Expected empty set to be empty", names[i])
		}
		if set.Size() != 0 {
			t.Errorf("%s: Expected size 0 for empty set", names[i])
		}
		if set.Contains(1) {
			t.Errorf("%s: Empty set should not contain any element", names[i])
		}
		if set.Remove(1) {
			t.Errorf("%s: Remove from empty set should return false", names[i])
		}
		slice := set.ToSlice()
		if len(slice) != 0 {
			t.Errorf("%s: ToSlice on empty set should return empty slice", names[i])
		}
	}
}

func TestSingleElement(t *testing.T) {
	hashSet := NewHashSet[string]()
	hashSet.Add("test")

	if hashSet.Size() != 1 {
		t.Error("Expected size 1")
	}
	if !hashSet.Contains("test") {
		t.Error("Expected to contain 'test'")
	}
	if hashSet.IsEmpty() {
		t.Error("Expected not to be empty")
	}

	slice := hashSet.ToSlice()
	if len(slice) != 1 || slice[0] != "test" {
		t.Error("ToSlice should return single element")
	}

	hashSet.Remove("test")
	if !hashSet.IsEmpty() {
		t.Error("Expected to be empty after removing only element")
	}
}
