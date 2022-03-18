package lib

import (
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SubStruct struct {
	C string `json:"c" url:"c"`
}

type ParamsStructExample struct {
	A string      `json:"a" url:"a" required:"true"`
	B string      `json:"-" url:"-"`
	C []SubStruct `json:"sub" url:"sub"`
	D []string    `json:"d" url:"d" required:"true"`
}

func Test_Params_ToJSON(t *testing.T) {
	p := Params{
		Params: ParamsStructExample{
			A: "The a value",
			B: "The b value",
			C: []SubStruct{{C: "the c value"}},
			D: []string{"hello"},
		},
	}

	reader, err := p.ToJSON()
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "{\"a\":\"The a value\",\"sub\":[{\"c\":\"the c value\"}],\"d\":[\"hello\"]}", string(b))
}

func Test_Params_ToJSON_Missing_Require(t *testing.T) {
	p := Params{
		Params: ParamsStructExample{
			A: "The a value",
			B: "The b value",
			C: []SubStruct{{C: "the c value"}},
		},
	}

	_, err := p.ToJSON()
	assert.Error(t, err, "")

	p = Params{
		Params: ParamsStructExample{
			B: "The b value",
			C: []SubStruct{{C: "the c value"}},
			D: []string{"hello"},
		},
	}

	_, err = p.ToJSON()
	assert.Error(t, err, "")
}

func Test_Params_ToValues(t *testing.T) {
	p := Params{
		Params: ParamsStructExample{
			A: "The a value",
			B: "The b value",
			C: []SubStruct{{C: "the c value"}},
			D: []string{"hello"},
		},
	}

	values, err := p.ToValues()
	assert.NoError(t, err)
	unescaped, _ := url.QueryUnescape(values.Encode())
	assert.Equal(t, "a=The a value&d[0]=hello&sub[0][c]=the c value", unescaped)
}
