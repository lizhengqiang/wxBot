package bot

// 这里来存储需要的数据
type Cache interface {
	Get(string) string
	Set(string, string) error
}


