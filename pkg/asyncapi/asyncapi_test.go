package asyncapi

import (
	"bytes"
	"fmt"
	"github.com/jensneuse/graphql-go-tools/pkg/astprinter"
	"github.com/stretchr/testify/require"
	"testing"
)

var asyncapiDoc = []byte(`
asyncapi: '2.0.0'
info:
  title: Streetlights API
  version: '1.0.0'
  description: |
    The Smartylighting Streetlights API allows you
    to remotely manage the city lights.
  license:
    name: Apache 2.0
    url: 'https://www.apache.org/licenses/LICENSE-2.0'
servers:
  mosquitto:
    url: mqtt://test.mosquitto.org
    protocol: mqtt
channels:
  light/measured:
    publish:
      summary: Inform about environmental lighting conditions for a particular streetlight.
      operationId: onLightMeasured
      message:
        name: LightMeasured
        payload:
          type: object
          properties:
            id:
              type: integer
              minimum: 0
              description: Id of the streetlight.
            lumens:
              type: integer
              minimum: 0
              description: Light intensity measured in lumens.
            sentAt:
              type: string
              format: date-time
              description: Date and time when the message was sent.`)

func TestAsyncAPI(t *testing.T) {
	doc, report := ImportAsyncAPIDocumentByte(asyncapiDoc)
	if report.HasErrors() {
		t.Fatal(report)
	}
	w := &bytes.Buffer{}
	err := astprinter.PrintIndent(&doc, nil, []byte(" "), w)
	require.NoError(t, err)
	fmt.Println(w.String())
}
