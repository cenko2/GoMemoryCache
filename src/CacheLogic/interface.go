package CacheLogic

type CacheHandlerInterface interface {
	Get(string) string
	Insert(string, *string, int)
	KeyExists(string) bool
}
