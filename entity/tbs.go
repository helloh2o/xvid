package entity

// GetTables 需要注册的表
func GetTables() []interface{} {
	return []interface{}{
		&AppUser{},
	}
}
