package templates

import (
	"io"
	"text/template"
)

const migrationJob = `apiVersion: batch/v1
kind: Job
metadata:
  name: uyuni-migration
spec:
  backoffLimit: 1
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: rsync-var-pgsql
        image: {{ .Image }}:{{ .Tag }}
        command: [ "/var/lib/uyuni-tools/migrate.sh" ]
        env:
          - name: SSH_AUTH_SOCK
            value: /tmp/ssh_auth_sock
        volumeMounts:
        {{- range $name, $volume := .Volumes }}
          - mountPath: {{ $volume.Path }}
            name: {{ $name }}
        {{- end }}
      volumes:
	  {{- range $name, $volume := .Volumes }}
        - name: {{ $name }}
		{{- if eq $volume.HostPath "" }}
          persistentVolumeClaim:
            claimName: {{ $name }}
		{{- else }}
          hostPath:
            path: {{ $volume.HostPath }}
	    {{- end }}
      {{- end }}
`

type Volume struct {
	HostPath string
	Path     string
}

type KubernetesMigrateJobTemplateData struct {
	Volumes map[string]Volume
	Image   string
	Tag     string
}

func (data KubernetesMigrateJobTemplateData) Render(wr io.Writer) error {
	t := template.Must(template.New("job").Parse(podmanMigrationScriptTemplate))
	return t.Execute(wr, data)
}
