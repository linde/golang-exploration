package restserver

import (
	"fmt"
	"io"
	"myapp/greeterserver"
	"myapp/grpcservice"
	"net/http"
	"testing"

	a "github.com/stretchr/testify/assert"
)

func TestRestGateway(t *testing.T) {

	assert := a.New(t)
	assert.NotNil(assert)

	// set up a net listenter so we can test gatewaying requests to it

	gs, err := grpcservice.NewServerFromPort(0)
	assert.NotNil(gs)
	assert.Nil(err)
	serverAssignedPort, portErr := gs.GetServicePort()
	assert.Nil(portErr)
	assert.Greater(serverAssignedPort, 0)

	helloServer := greeterserver.NewHelloServer()
	defer helloServer.Stop()
	go gs.Serve(helloServer)

	// now get our grpc address and set up a rest gateway proxying to it.
	//  we use a channel to receive it's address
	gsTCPAddr, addressErr := gs.GetServiceTCPAddr()
	assert.Nil(addressErr)
	rgw := NewRestGateway(0, gsTCPAddr)
	assert.NotNil(rgw)
	defer rgw.Close()

	gwAddr := rgw.GetRestGatewayAddr()
	go rgw.Serve()

	// ok, we're set up to test the restgawy

	tests := []struct {
		nameInput  string
		timesInput int
		pauseInput int
		respCode   int
	}{
		{"dolly", 1, 0, http.StatusOK},
		{"dolly", 2, 0, http.StatusOK},
		{"negativeTimes", -3, 0, http.StatusBadRequest},
		{"negativeWait", 1, -1, http.StatusBadRequest},
	}

	for idx, test := range tests {

		testName := fmt.Sprintf("TestRestGateway(idx:%d){name:%s,times:%d,pause:%d}",
			idx, test.nameInput, test.timesInput, test.pauseInput)

		t.Run(testName, func(tt *testing.T) {

			assert := a.New(tt)

			url := fmt.Sprintf("http://%s/v1/helloservice/sayhello?name=%s&times=%d&pause=%d",
				gwAddr, test.nameInput, test.timesInput, test.pauseInput)
			resp, httpErr := http.Get(url)
			assert.Nil(httpErr)
			assert.NotNil(resp)
			assert.Equal(test.respCode, resp.StatusCode)
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.Nil(err)

				bodyStr := string(body)
				assert.Contains(bodyStr, test.nameInput)
				assert.Contains(bodyStr, fmt.Sprintf("%d of %d", test.timesInput, test.timesInput))
			}
		})
	}
}
