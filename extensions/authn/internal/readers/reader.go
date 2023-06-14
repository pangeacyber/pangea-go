package readers

type Reader interface {
	Next() (map[string]interface{}, error)
}
