package iswbx

import (
	"bytes"
	"text/template"

	"github.com/sirupsen/logrus"
)

func alarmToMessage(wh map[string]interface{}) string {
	var tmpl *template.Template
	if operation, ok := wh["Operation"]; ok {
		switch operation {
		case "Modified":
			tmpl = alarmModifiedTmpl
		case "Created":
			tmpl = alarmCreatedTmpl
		default:
			logrus.Errorf("Unsupported event operation: %v", operation)
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

	logrus.Error("Webhook missing operation")
	return ""
}

func webhookToMessage(wh map[string]interface{}) string {
	if eventObjectType, ok := wh["EventObjectType"]; ok {
		switch eventObjectType {
		case "cond.Alarm":
			return alarmToMessage(wh)
		default:
			logrus.WithField("wh", wh).Errorf("Unsupported event object type: %v", eventObjectType)
			return ""
		}
	}

	logrus.Error("Event has no object type")
	return ""

}

var alarmModifiedTmpl = template.Must(template.New("name").Parse(`
## Intersight Alarm Updated

**Severity:** {{ .Event.Severity }} (Original Severity: {{ .Event.OrigSeverity }})

**Affected Object:** {{ .Event.AffectedMoDisplayName }} ({{ .Event.AffectedMoType }})

**Message:** {{ .Event.Name }}: {{ .Event.Description }}

**Creation Time:** {{ .Event.CreateTime }}

**Last Transition Time:** {{ .Event.LastTransitionTime }}`))
var alarmCreatedTmpl = template.Must(template.New("name").Parse(`
## Intersight Alarm Created

**Severity:** {{ .Event.Severity }}

**Affected Object:** {{ .Event.AffectedMoDisplayName }} ({{ .Event.AffectedMoType }})

**Message:** {{ .Event.Name }}: {{ .Event.Description }}

**Creation Time:** {{ .Event.CreateTime }}

**Last Transition Time:** {{ .Event.LastTransitionTime }}`))
