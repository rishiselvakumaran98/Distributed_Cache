package cache

import "time"

type CacheSet interface {
	// Has: function checks if the Given key exists
	Has([]byte) bool 
	// Get: function returns value if found or error 
	Get([]byte) ([]byte, error)
	// Set: function sets the key and value in cache
	Set([]byte, []byte, time.Duration ) error
	// Delete: function deletes the given key value and returns error if not found
	Delete([]byte) error

	// TODO: Error handling in project

}