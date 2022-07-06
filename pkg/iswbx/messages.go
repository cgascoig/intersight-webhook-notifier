package iswbx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
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
	var operation string
	if op, ok := wh["Operation"]; ok {
		if op, ok := op.(string); ok {
			operation = op
		} else {
			logrus.Error("Webhook missing operation")
			return ""
		}
	} else {
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
	// case "workflow.WorkflowInfo":
	// 	switch operation {
	// 	case "Modified":
	// 		tmpl = workflowModifiedTmpl
	// 	case "Created":
	// 		tmpl = workflowCreatedTmpl
	// 	default:
	// 		logrus.Errorf("Unsupported event operation: %v", operation)
	// 		return ""
	// 	}
	case "":
		logrus.Info("Received Intersight PING Webhook")
		if subscription, ok := wh["Subscription"]; ok {
			if subscription, ok := subscription.(map[string]interface{}); ok {
				if m, ok := subscription["Moid"]; ok {
					if m, ok := m.(string); ok {
						return fmt.Sprintf("Intersight webhooks are now being received for subscription [%s](https://www.intersight.com/an/settings/webhooks/%s/edit/)", m, m)
					}
				}
			}
		}

		return "Intersight webhooks are now being received"
	default:
		logrus.WithField("wh", wh).Warnf("Unsupported event object type: %v", eventObjectType)

		event, ok := wh["Event"]
		if !ok {
			logrus.Errorf("Unsupported event object type and with no Event")
			return ""
		}

		eventJSON, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			logrus.Errorf("Error marshalling JSON for generic/unsupported event")
			return ""
		}

		msgFormat := `
## Intersight %s %s

An Intersight event was received with an event type that I don't support yet, but here is the %s object:

` + "```\n%s\n```"
		msg := fmt.Sprintf(msgFormat, eventObjectType, operation, strings.ToLower(operation), eventJSON)
		return msg
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
