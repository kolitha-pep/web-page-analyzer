package utils

var cacheMap = make(map[string]interface{})

// SetCache sets a value in the cache with a given key.
func SetCache(key string, value interface{}) {
	cacheMap[key] = value
}

// GetCache retrieves a value from the cache by its key.
func GetCache(key string) (interface{}, bool) {
	value, exists := cacheMap[key]
	return value, exists
}
