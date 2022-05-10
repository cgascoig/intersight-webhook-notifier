package iswbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookToMessage(t *testing.T) {
	tests := []struct {
		in  map[string]interface{}
		out string
	}{
		{
			in: map[string]interface{}{
				"ClassId":         "mo.WebhookResult",
				"ObjectType":      "mo.WebhookResult",
				"EventObjectType": "cond.Alarm",
				"Operation":       "Modified",
				"Event": map[string]interface{}{
					"Acknowledge":           "None",
					"AcknowledgeBy":         "",
					"AcknowledgeTime":       "0001-01-01T00:00:00Z",
					"AffectedMoDisplayName": "Hyperflex/sys/chassis-1/blade-5/mgmt/actual-mount-list/actual-mount-entry-1",
					"AffectedMoId":          "60b7211076752d313426fbbc",
					"AffectedMoType":        "management.Controller",
					"AffectedObject":        "60b720f96f72612d31e9b344/sys/chassis-1/blade-5/mgmt/actual-mount-list/actual-mount-entry-1/fault-F1404",
					"AncestorMoId":          "60b7210776752d313426f534",
					"AncestorMoType":        "compute.Blade",
					"ClassId":               "cond.Alarm",
					"Code":                  "F1404",
					"CreateTime":            "2022-03-08T13:36:48.491Z",
					"CreationTime":          "2022-03-08T13:35:51.676Z",
					"Description":           "Server 1/5 (service profile: org-root/org-grscarle/ls-Linux-iSCSI) vmedia mapping CentOS-7-ISO has failed.",
					"DomainGroupMoid":       "5b25418d7a7662743465cf72",
					"LastTransitionTime":    "2022-05-09T05:55:17.53Z",
					"ModTime":               "2022-05-09T05:56:25.01Z",
					"Moid":                  "62275bf065696e2d33b4ba3c",
					"MsAffectedObject":      "sys/chassis-1/blade-5/mgmt/actual-mount-list/actual-mount-entry-1",
					"Name":                  "UCS-F1404",
					"ObjectType":            "cond.Alarm",
					"OrigSeverity":          "Critical",
					"Severity":              "Critical",
				},
			},
			out: `
## Intersight Alarm Updated

**Severity:** Critical (Original Severity: Critical)

**Affected Object:** Hyperflex/sys/chassis-1/blade-5/mgmt/actual-mount-list/actual-mount-entry-1 (management.Controller)

**Message:** UCS-F1404: Server 1/5 (service profile: org-root/org-grscarle/ls-Linux-iSCSI) vmedia mapping CentOS-7-ISO has failed.

**Creation Time:** 2022-03-08T13:36:48.491Z

**Last Transition Time:** 2022-05-09T05:55:17.53Z`,
		},
		{
			in: map[string]interface{}{
				"ClassId":         "mo.WebhookResult",
				"ObjectType":      "mo.WebhookResult",
				"EventObjectType": "cond.Alarm",
				"Operation":       "Created",
				"Event": map[string]interface{}{
					"Acknowledge":           "None",
					"AcknowledgeBy":         "",
					"AcknowledgeTime":       "0001-01-01T00:00:00Z",
					"AffectedMoDisplayName": "Hyperflex/sys/chassis-1/blade-5/mgmt/actual-mount-list/actual-mount-entry-1",
					"AffectedMoId":          "60b7211076752d313426fbbc",
					"AffectedMoType":        "management.Controller",
					"AffectedObject":        "60b720f96f72612d31e9b344/sys/chassis-1/blade-5/mgmt/actual-mount-list/actual-mount-entry-1/fault-F1404",
					"AncestorMoId":          "60b7210776752d313426f534",
					"AncestorMoType":        "compute.Blade",
					"ClassId":               "cond.Alarm",
					"Code":                  "F1404",
					"CreateTime":            "2022-03-08T13:36:48.491Z",
					"CreationTime":          "2022-03-08T13:35:51.676Z",
					"Description":           "Server 1/5 (service profile: org-root/org-grscarle/ls-Linux-iSCSI) vmedia mapping CentOS-7-ISO has failed.",
					"DomainGroupMoid":       "5b25418d7a7662743465cf72",
					"LastTransitionTime":    "2022-05-09T05:55:17.53Z",
					"ModTime":               "2022-05-09T05:56:25.01Z",
					"Moid":                  "62275bf065696e2d33b4ba3c",
					"MsAffectedObject":      "sys/chassis-1/blade-5/mgmt/actual-mount-list/actual-mount-entry-1",
					"Name":                  "UCS-F1404",
					"ObjectType":            "cond.Alarm",
					"OrigSeverity":          "Critical",
					"Severity":              "Critical",
				},
			},
			out: `
## Intersight Alarm Created

**Severity:** Critical

**Affected Object:** Hyperflex/sys/chassis-1/blade-5/mgmt/actual-mount-list/actual-mount-entry-1 (management.Controller)

**Message:** UCS-F1404: Server 1/5 (service profile: org-root/org-grscarle/ls-Linux-iSCSI) vmedia mapping CentOS-7-ISO has failed.

**Creation Time:** 2022-03-08T13:36:48.491Z

**Last Transition Time:** 2022-05-09T05:55:17.53Z`,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.out, webhookToMessage(test.in))
	}
}
