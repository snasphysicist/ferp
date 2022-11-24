package mapper

import "testing"

func TestMapWithTypeForwardUnchangedOnlyDeserialisesIntoPassthrough(t *testing.T) {
	m := map[string]string{"type": "forward-unchanged"}
	err := (&Passthrough{}).From(m)
	if err != nil {
		t.Errorf("Failed to deserialise %+v into Passthrough: %s", m, err)
	}
}

func TestMapWithFieldsOtherThanTypeDoesNotDeserialiseIntoPassthrough(t *testing.T) {
	m := map[string]string{"type": "forward-unchanged", "foo": "bar"}
	err := (&Passthrough{}).From(m)
	if err == nil {
		t.Errorf("Deserialised %+v into Passthrough, but should not be possible", m)
	}
}

func TestMapWithTypeOtherThanForwardUnchangedDoesNotDeserialiseIntoPassthrough(t *testing.T) {
	m := map[string]string{"type": "something-else"}
	err := (&Passthrough{}).From(m)
	if err == nil {
		t.Errorf("Deserialised %+v into Passthrough, but should not be possible", m)
	}
}
