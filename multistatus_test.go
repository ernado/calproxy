package main

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiStatus(t *testing.T) {
	raw := `<?xml version='1.0' encoding='utf-8'?>
<ns0:multistatus xmlns:ns0="DAV:">
    <ns0:response>
        <ns0:href>/principals/example.com/i.ivanov/calendars/edcd5ce4-ca1a-49f4-b705-98a56ea3dfec/</ns0:href>
        <ns0:propstat>
            <ns0:prop>
                <ns0:getetag>8d25c704-0ef3-4a66-862f-150ea60cb68d</ns0:getetag>
            </ns0:prop>
            <ns0:status>HTTP/1.1 200 OK</ns0:status>
        </ns0:propstat>
    </ns0:response>
    <ns0:response>
        <ns0:href>
            /principals/example.com/i.ivanov/calendars/edcd5ce4-ca1a-49f4-b705-98a56ea3dfec/b225a668-052c-4b05-bbe8-9593c7f9cc8e.ics
        </ns0:href>
        <ns0:propstat>
            <ns0:prop>
                <ns0:getetag>065845df-af03-4aaa-a473-2372404eef5e</ns0:getetag>
            </ns0:prop>
            <ns0:status>HTTP/1.1 200 OK</ns0:status>
        </ns0:propstat>
    </ns0:response>
</ns0:multistatus>
`

	status, err := DecodeMultiStatus([]byte(raw))
	require.NoError(t, err)

	out, err := xml.MarshalIndent(status, "", "  ")
	require.NoError(t, err)

	t.Log(string(out))
}
