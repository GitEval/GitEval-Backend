package tool

// SafeString 检查 string 指针是否为 nil，如果为 nil 则返回空字符串
func SafeString(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}

// SafeInt 检查 int 指针是否为 nil，如果为 nil 则返回 0
func SafeInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

// SafeInt64 检查 int64 指针是否为 nil，如果为 nil 则返回 0
func SafeInt64(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}

// SafeInt32 检查 int32 指针是否为 nil，如果为 nil 则返回 0
func SafeInt32(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}

// SafeFloat64 检查 float64 指针是否为 nil，如果为 nil 则返回 0
func SafeFloat64(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0.0
}

// SafeBool 检查 bool 指针是否为 nil，如果为 nil 则返回 false
func SafeBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

// SafeUint64 检查 uint64 指针是否为 nil，如果为 nil 则返回 0
func SafeUint64(u *uint64) uint64 {
	if u != nil {
		return *u
	}
	return 0
}

// SafeSlice 检查 slice 指针是否为 nil，如果为 nil 则返回空切片
func SafeSlice[T any](s *[]T) []T {
	if s != nil {
		return *s
	}
	return []T{}
}

// SafeMap 检查 map 指针是否为 nil，如果为 nil 则返回空 map
func SafeMap[K comparable, V any](m *map[K]V) map[K]V {
	if m != nil {
		return *m
	}
	return make(map[K]V)
}

// SafePointer 检查指针是否为 nil，如果为 nil 则返回 nil
func SafePointer[T any](p *T) *T {
	if p != nil {
		return p
	}
	return nil
}
