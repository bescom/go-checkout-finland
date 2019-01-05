package checkout

import (
	"testing"
)

type teststruct struct {
	secret string
	data   string
	result string
}

var tests = []teststruct{
	{"", "", ""},
	{"secret", "", ""},
	{"", "data", ""},
	{"secret", "data", "1b2c16b75bd2a870c114153ccda5bcfca63314bc722fa160d690de133ccbb9db"},
}

func TestGenerateHMAC(t *testing.T) {

	for _, test := range tests {

		hmac, _ := GenerateHMAC(test.secret, test.data)

		if hmac != test.result {

			t.Error(
				"HMAC("+test.secret+","+test.data+")",
				"expected", test.result,
				",got", hmac,
			)

		}

	}

}
