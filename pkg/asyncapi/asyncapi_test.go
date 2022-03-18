package asyncapi

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/jensneuse/graphql-go-tools/pkg/astprinter"
	"github.com/stretchr/testify/require"
)

func TestAsyncAPI_Streetlights(t *testing.T) {
	asyncapiDoc, err := ioutil.ReadFile("./fixtures/streetlights.yaml")
	require.NoError(t, err)

	doc, report := ImportAsyncAPIDocumentByte(asyncapiDoc)
	if report.HasErrors() {
		t.Fatal(report)
	}
	w := &bytes.Buffer{}
	err = astprinter.PrintIndent(&doc, nil, []byte(" "), w)
	require.NoError(t, err)
	result := w.Bytes()

	fixture, err := ioutil.ReadFile("./fixtures/streetlights.graphql")
	require.NoError(t, err)

	require.Equal(t, fixture, result)
}

func TestAsyncAPI_Streetlights_Kafka(t *testing.T) {
	asyncapiDoc, err := ioutil.ReadFile("./fixtures/streetlights-kafka.yaml")
	require.NoError(t, err)

	doc, report := ImportAsyncAPIDocumentByte(asyncapiDoc)
	if report.HasErrors() {
		t.Fatal(report)
	}
	w := &bytes.Buffer{}
	err = astprinter.PrintIndent(&doc, nil, []byte("  "), w)
	require.NoError(t, err)
	result := w.Bytes()

	fixture, err := ioutil.ReadFile("./fixtures/streetlights-kafka.graphql")
	require.NoError(t, err)

	require.Equal(t, string(fixture), string(result))
}

func TestAsyncAPI_EmailService(t *testing.T) {
	asyncapiDoc, err := ioutil.ReadFile("./fixtures/email-service.yaml")
	require.NoError(t, err)

	doc, report := ImportAsyncAPIDocumentByte(asyncapiDoc)
	if report.HasErrors() {
		t.Fatal(report)
	}
	w := &bytes.Buffer{}
	err = astprinter.PrintIndent(&doc, nil, []byte("  "), w)
	require.NoError(t, err)
	result := w.Bytes()

	fixture, err := ioutil.ReadFile("./fixtures/email-service.graphql")
	require.NoError(t, err)
	require.Equal(t, string(fixture), string(result))
}
