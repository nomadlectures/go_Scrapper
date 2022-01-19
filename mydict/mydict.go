package mydict

import "errors"

// Dictionary type
type Dictionary map[string]string

var (
	errNotFound = errors.New("Not Found")

	errWordExists = errors.New("Already exists")
	errCantUpdate = errors.New("Cant update non-existing word")
)

// Search for a word
func (d Dictionary) Search(word string) (string, error) {

	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFound
}

// Add the new element
func (d Dictionary) Add(word string, def string) error {

	_, err := d.Search(word)

	if err == errNotFound {

		d[word] = def
	} else if err == nil {
		return errWordExists
	}

	return nil
}

// Update a element
func (d Dictionary) Update(word string, def string) error {
	_, err := d.Search(word)

	switch err {
	case nil:
		d[word] = def
	case errNotFound:
		return errCantUpdate
	}

	return nil
}

// Delete a word
func (d Dictionary) Delete(word string) {

	delete(d, word)
}
