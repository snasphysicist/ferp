package mapper

import (
	"fmt"
	"log"
	"strings"
)

// RemovePrefix is a path mapper which removes a prefix from the path
// before forwarding the request to the downstream
type RemovePrefix struct {
	prefix string
}

// NewRemovePrefix creates a new RemovePrefix which removes the given prefix
func NewRemovePrefix(prefix string) RemovePrefix {
	return RemovePrefix{prefix: prefix}
}

// Map implements pathRewriter for RemovePrefix
func (m RemovePrefix) Map(from string) string {
	t := strings.TrimPrefix(from, m.prefix)
	log.Printf("Removing prefix %s: rewriting %s to %s", m.prefix, from, t)
	return t
}

// From deserialises a string-string map into a RemovePrefix
func (m *RemovePrefix) From(c map[string]string) error {
	if len(c) != 2 {
		return fmt.Errorf(
			"configuration %+v has %d fields, RemovePrefix requires exactly 2 (type == remove-prefix, prefix)",
			c, len(c))
	}
	t, tOK := c["type"]
	if !tOK {
		return fmt.Errorf(
			"configuration %+v has no type field, required for RemovePrefix", c)
	}
	if t != "remove-prefix" {
		return fmt.Errorf("configuration %+v has type %s, RemovePrefix requires remove-prefix", c, t)
	}
	p, pOK := c["prefix"]
	if !pOK {
		return fmt.Errorf(
			"configuration %+v has no prefix field, required for RemovePrefix", c)
	}
	m.prefix = p
	return nil
}
