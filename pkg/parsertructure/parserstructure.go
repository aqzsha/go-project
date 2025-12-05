package parsertructure

import "encoding/json"

func Parse[T any](src interface{}, dest *T) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, dest)
}
