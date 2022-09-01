package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	json1 := `{"a":{"b":"c"}}`
	data, err := Decode([]byte(json1))
	assert.NoError(t, err)

	ab := data.GetString("a.b")
	assert.Equal(t, "c", ab)

	json2, err := data.Json()
	assert.NoError(t, err)
	assert.Equal(t, `{"a":{"b":"c"}}`, string(json2))

	err = data.Set("a.b", "c1")
	assert.NoError(t, err)

	json3, err := data.Json()
	assert.NoError(t, err)
	assert.Equal(t, `{"a":{"b":"c1"}}`, string(json3))
}

func TestEncode(t *testing.T) {
	type foo struct {
		A struct {
			B string `json:"b"`
		} `json:"a"`
	}

	var f foo
	f.A.B = "c"

	data, err := Encode(f)
	assert.NoError(t, err)

	ab := data.GetString("a.b")
	assert.Equal(t, "c", ab)

	json2, err := data.Json()
	assert.NoError(t, err)
	assert.Equal(t, `{"a":{"b":"c"}}`, string(json2))

	err = data.Set("a.b", "c1")
	assert.NoError(t, err)

	json3, err := data.Json()
	assert.NoError(t, err)
	assert.Equal(t, `{"a":{"b":"c1"}}`, string(json3))

	var want, got foo
	want.A.B = "c1"
	err = data.Read("", &got)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
