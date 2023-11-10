package templates

import (
	"io"
	"text/template"
)

const networkTemplate = `
[Network]
Options=
NetworkName={{ .Network }}
`

type PodmanNetworkTemplateData struct {
	Network    string
}

func (data PodmanNetworkTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("network").Parse(networkTemplate))
	return t.Execute(wr, data)
}
