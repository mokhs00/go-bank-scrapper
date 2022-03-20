package dict

import "errors"

// Dictionary type
type Dictionary map[string]string

var errNotFoundElement = errors.New("not found element")

func (dictionary Dictionary) Search(word string) (string, error) {
	value, exists := dictionary[word]

	if exists {
		return value, nil
	}

	return "", errNotFoundElement
}
