package importio

type Reader interface {
	Next() (map[string]interface{}, error)
}
