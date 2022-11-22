package mapper

import (
	"fmt"
	"log"
)

// Passthrough is a path mapper which does not modify the path at all for the downstream
type Passthrough struct{}

// Map implements pathRewriter for Passthrough
func (Passthrough) Map(from string) string {
	log.Printf("Rewriting %s to %s", from, from)
	return from
}

// From deserialises a string-string map into a Passthrough
func (m *Passthrough) From(c map[string]string) error {
	if len(c) != 1 {
		return fmt.Errorf(
			"configuration %+v has %d fields, Passthrough requires exactly 1 (type == forward-unchanged)",
			c, len(c))
	}
	t, tOK := c["type"]
	if !tOK {
		return fmt.Errorf(
			"configuration %+v has no type field, required for Passthrough", c)
	}
	if t != "forward-unchanged" {
		return fmt.Errorf("configuration %+v has type %s, Passthrough requires forward-unchanged", c, t)
	}
	return nil
}
