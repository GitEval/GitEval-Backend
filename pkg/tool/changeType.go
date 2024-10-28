package tool

// 辅助函数，检查指针类型字段的值
func SafeString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
}

func SafeInt32(ptr *int) int32 {
	if ptr != nil {
		return int32(*ptr)
	}
	return 0
}
