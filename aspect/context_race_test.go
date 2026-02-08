// Package aspect - context_race_test checks for race conditions within a single Context instance,
// primarily focusing on the Metadata map accessed by different advice functions.
package aspect

import (
	"fmt"
	"sync"
	"testing"
)

func TestContextMetadataRace_SimulatedConcurrentAccess(t *testing.T) {
	// Note: The standard execution engine runs advice sequentially based on type/priority.
	// This test simulates a scenario where multiple goroutines *could* access the same
	// Context's Metadata map simultaneously, testing its internal safety if used this way.
	// This is more about testing the *conceptual* safety of the Metadata map itself.

	c := NewContext("SimulatedFunc")
	key := "shared_key"
	value := "test_value"
	generateAlphabeticalKeyInc := func(index int, key, topic string) string {
		return fmt.Sprintf("%s_%s_%s", key, topic, string(rune('A'+index%10)))
	}

	var wg sync.WaitGroup

	// Goroutine writing to metadata
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			c.SetMetadataVal(generateAlphabeticalKeyInc(i, key, "write"), value)
		}
	}()

	// Goroutine reading from metadata
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_, _ = c.GetMetadataVal(generateAlphabeticalKeyInc(i, key, "write")) // Read
		}
	}()

	// Another goroutine writing
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			c.SetMetadataVal(generateAlphabeticalKeyInc(i, key, "write2"), value) // Write
		}
	}()

	wg.Wait()
	// If the map is not safe for concurrent read/write, 'go test -race' will detect it.
}
