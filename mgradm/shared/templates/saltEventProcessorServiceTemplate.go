package templates

import (
	"io"
	"text/template"
)

// SaltEventProcessorServiceTemplateData represents the data for salt event processor service.
const saltEventProcessorServiceTemplate = `
[Unit]
Description=Uyuni Event Processor Container
Wants=network.target
After=network-online.target

[Service]
Environment=PODMAN_SYSTEMD_UNIT=%n
Restart=on-failure
ExecStartPre=/bin/rm -f %t/uyuni-salt-event-processor-%i.pid %t/%n.ctr-id
ExecStartPre=/usr/bin/podman rm --ignore --force -t 10 {{ .NamePrefix }}-salt-event-processor-%i
ExecStart=/bin/sh -c '/usr/bin/podman run \
    --conmon-pidfile %t/uyuni-salt-event-processor-%i.pid \
    --cidfile=%t/%n-%i.ctr-id \
    --cgroups=no-conmon \
    --sdnotify=conmon \
    -d \
    -e db_name=${UYUNI_DB_NAME} \
	-e db_port=${UYUNI_DB_PORT} \
    -e db_host=${UYUNI_DB_HOST} \
	-e db_backend=postgresql \
    --secret={{ .DBUserSecret }},type=env,target=db_user \
    --secret={{ .DBPassSecret }},type=env,target=db_password \
    --replace \
    --name {{ .NamePrefix }}-salt-event-processor-%i \
    --hostname {{ .NamePrefix }}-server-event-processor-%i.mgr.internal \
    --network {{ .Network }} \
    ${UYUNI_EVENT_PROCESSOR_IMAGE}'
ExecStop=/usr/bin/podman stop --ignore -t 10 --cidfile=%t/%n-%i.ctr-id
ExecStopPost=/usr/bin/podman rm -f --ignore -t 10 --cidfile=%t/%n-%i.ctr-id
PIDFile=%t/uyuni-salt-event-processor-%i.pid
TimeoutStopSec=60
TimeoutStartSec=60
Type=forking

[Install]
WantedBy=multi-user.target default.target
`

type EventProcessorServiceTemplateData struct {
	NamePrefix   string // "uyuni"
	Image        string
	Network      string // "uyuni-server"
	DBUserSecret string // "uyuni-db-user"
	DBPassSecret string // "uyuni-db-pass"
	DBBackend    string // "postgresql"
}

// Render will create the systemd configuration file.
func (data EventProcessorServiceTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("service").Parse(saltEventProcessorServiceTemplate))
	return t.Execute(wr, data)
}
