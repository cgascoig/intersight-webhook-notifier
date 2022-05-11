package iswbx

import (
	"bytes"
	"text/template"

	"github.com/sirupsen/logrus"
)

func webhookToMessage(wh map[string]interface{}) string {
	var tmpl *template.Template

	eventObjectType, ok := wh["EventObjectType"]
	if !ok {
		logrus.Error("Event has no object type")
		return ""
	}

	operation, ok := wh["Operation"]
	if !ok {
		logrus.Error("Webhook missing operation")
		return ""
	}

	switch eventObjectType {
	case "cond.Alarm":
		// return alarmToMessage(wh)
		switch operation {
		case "Modified":
			tmpl = alarmModifiedTmpl
		case "Created":
			tmpl = alarmCreatedTmpl
		default:
			logrus.Errorf("Unsupported event operation: %v", operation)
			return ""
		}

	case "workflow.WorkflowInfo":
		switch operation {
		case "Modified":
			tmpl = workflowModifiedTmpl
		case "Created":
			tmpl = workflowCreatedTmpl
		default:
			logrus.Errorf("Unsupported event operation: %v", operation)
			return ""
		}
	default:
		logrus.WithField("wh", wh).Errorf("Unsupported event object type: %v", eventObjectType)
		return ""
	}

	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, wh)
	if err != nil {
		logrus.Errorf("Template execution failed: %v", err)
		return ""
	}
	return buf.String()
}

var alarmModifiedTmpl = template.Must(template.New("name").Parse(`
## Intersight Alarm Updated

**Severity:** {{ .Event.Severity }} (Original Severity: {{ .Event.OrigSeverity }})

**Affected Object:** {{ .Event.AffectedMoDisplayName }} ({{ .Event.AffectedMoType }})

**Message:** [{{ .Event.Name }}](https://www.intersight.com/an/cond/alarms/active/?Moid={{ .Event.Moid }}): {{ .Event.Description }}

**Creation Time:** {{ .Event.CreateTime }}

**Last Transition Time:** {{ .Event.LastTransitionTime }}`))
var alarmCreatedTmpl = template.Must(template.New("name").Parse(`
## Intersight Alarm Created

**Severity:** {{ .Event.Severity }}

**Affected Object:** {{ .Event.AffectedMoDisplayName }} ({{ .Event.AffectedMoType }})

**Message:** [{{ .Event.Name }}](https://www.intersight.com/an/cond/alarms/active/?Moid={{ .Event.Moid }}): {{ .Event.Description }}

**Creation Time:** {{ .Event.CreateTime }}

**Last Transition Time:** {{ .Event.LastTransitionTime }}`))

var workflowCreatedTmpl = template.Must(template.New("name").Parse(`{{ if not .Event.Internal }}
## Intersight Workflow Started

**Name:** [{{ .Event.Name }}](https://www.intersight.com/an/workflow/workflow-infos/{{ .Event.Moid }})

**Email:** {{ .Event.Email }}

**Status:** {{ .Event.Status }}

**Start Time:** {{ .Event.StartTime }}{{ end }}`))
var workflowModifiedTmpl = template.Must(template.New("name").Parse(`{{ if not .Event.Internal }}{{ if ne .Event.Status "RUNNING" }}
## Intersight Workflow {{ if eq .Event.Status "COMPLETED" }}Completed{{ else }}Updated{{ end }}

**Name:** [{{ .Event.Name }}](https://www.intersight.com/an/workflow/workflow-infos/{{ .Event.Moid }})

**Email:** {{ .Event.Email }}

**Status:** {{ .Event.Status }}

**Start Time:** {{ .Event.StartTime }}{{ end }}{{ else }}{{ if eq .Event.Status "TERMINATED" }}
## Intersight Workflow Terminated

**Name:** [{{ .Event.Name }}](https://www.intersight.com/an/workflow/workflow-infos/{{ .Event.Moid }})

**Email:** {{ .Event.Email }}

**Status:** {{ .Event.Status }}

**Start Time:** {{ .Event.StartTime }}{{ end }}{{ end }}`))
