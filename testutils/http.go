package testutils

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func DoHttpTest(t *testing.T, url string, expectedRespCode int, assertionStrings []string) {

	assert := assert.New(t)

	resp, httpErr := http.Get(url)
	assert.Nil(httpErr)
	assert.NotNil(resp)
	defer resp.Body.Close()

	assert.Equal(expectedRespCode, resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		assert.Nil(err)
		bodyStr := string(body)

		for _, assertStr := range assertionStrings {
			assert.Contains(bodyStr, assertStr)
		}
	}

}
