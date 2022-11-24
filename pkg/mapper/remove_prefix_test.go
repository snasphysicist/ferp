package mapper

import "testing"

func TestMapWithTypeRemovePrefixAndPrefixKeyDeserialisesIntoRemovePrefix(t *testing.T) {
	m := map[string]string{"type": "remove-prefix", "prefix": "foo"}
	err := (&RemovePrefix{}).From(m)
	if err != nil {
		t.Errorf("Failed to deserialise %+v into RemovePrefix: %s", m, err)
	}
}

func TestMapWithNoPrefixKeyDoesNotDeserialiseIntoRemovePrefix(t *testing.T) {
	m := map[string]string{"type": "remove-prefix", "fixpre": "foo"}
	err := (&RemovePrefix{}).From(m)
	if err == nil {
		t.Errorf("Deserialised %+v into RemovePrefix, but should not be possible", m)
	}
}

func TestMapWithTypeNotRemovePrefixDoesNotDeserialiseIntoRemovePrefix(t *testing.T) {
	m := map[string]string{"type": "something-else", "prefix": "foo"}
	err := (&RemovePrefix{}).From(m)
	if err == nil {
		t.Errorf("Deserialised %+v into RemovePrefix, but should not be possible", m)
	}
}
