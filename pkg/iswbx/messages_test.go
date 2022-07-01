package iswbx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookToMessage(t *testing.T) {
	tests := []struct {
		msg string
		in  map[string]interface{}
		out string
	}{
		{
			msg: "normal alarm modified message",
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

**Message:** [UCS-F1404](https://www.intersight.com/an/cond/alarms/active/?Moid=62275bf065696e2d33b4ba3c): Server 1/5 (service profile: org-root/org-grscarle/ls-Linux-iSCSI) vmedia mapping CentOS-7-ISO has failed.

**Creation Time:** 2022-03-08T13:36:48.491Z

**Last Transition Time:** 2022-05-09T05:55:17.53Z`,
		},
		{
			msg: "normal alarm created message",
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

**Message:** [UCS-F1404](https://www.intersight.com/an/cond/alarms/active/?Moid=62275bf065696e2d33b4ba3c): Server 1/5 (service profile: org-root/org-grscarle/ls-Linux-iSCSI) vmedia mapping CentOS-7-ISO has failed.

**Creation Time:** 2022-03-08T13:36:48.491Z

**Last Transition Time:** 2022-05-09T05:55:17.53Z`,
		},

		{
			msg: "normal workflow created message",
			in: map[string]interface{}{
				"ClassId":         "mo.WebhookResult",
				"ObjectType":      "mo.WebhookResult",
				"EventObjectType": "workflow.WorkflowInfo",
				"Operation":       "Created",
				"Event": map[string]interface{}{
					"ClassId":        "workflow.WorkflowInfo",
					"Email":          "email@cisco.com",
					"EndTime":        "2022-05-10T00:09:39.362Z",
					"StartTime":      "2022-05-10T00:07:15.844Z",
					"Internal":       false,
					"Moid":           "6279acb3696f6e2d31c5fcdd",
					"Name":           "RemoveKubernetesClusterProfileResources",
					"ObjectType":     "workflow.WorkflowInfo",
					"Parent":         nil,
					"ParentTaskInfo": nil,
					"Status":         "RUNNING",
				},
			},
			out: `
## Intersight Workflow Started

**Name:** [RemoveKubernetesClusterProfileResources](https://www.intersight.com/an/workflow/workflow-infos/6279acb3696f6e2d31c5fcdd)

**Email:** email@cisco.com

**Status:** RUNNING

**Start Time:** 2022-05-10T00:07:15.844Z`,
		},
		{
			msg: "internal workflows should generate empty message",
			in: map[string]interface{}{
				"ClassId":         "mo.WebhookResult",
				"ObjectType":      "mo.WebhookResult",
				"EventObjectType": "workflow.WorkflowInfo",
				"Operation":       "Created",
				"Event": map[string]interface{}{
					"ClassId":        "workflow.WorkflowInfo",
					"Email":          "email@cisco.com",
					"EndTime":        "2022-05-10T00:09:39.362Z",
					"StartTime":      "2022-05-10T00:07:15.844Z",
					"Internal":       true,
					"Moid":           "6279acb3696f6e2d31c5fcdd",
					"Name":           "RemoveKubernetesClusterProfileResources",
					"ObjectType":     "workflow.WorkflowInfo",
					"Parent":         nil,
					"ParentTaskInfo": nil,
					"Status":         "RUNNING",
				},
			},
			out: ``,
		},
		{
			msg: "updates to running workflows generate no message",
			in: map[string]interface{}{
				"ClassId":         "mo.WebhookResult",
				"ObjectType":      "mo.WebhookResult",
				"EventObjectType": "workflow.WorkflowInfo",
				"Operation":       "Modified",
				"Event": map[string]interface{}{
					"ClassId":        "workflow.WorkflowInfo",
					"Email":          "email@cisco.com",
					"EndTime":        "2022-05-10T00:09:39.362Z",
					"StartTime":      "2022-05-10T00:07:15.844Z",
					"Internal":       false,
					"Moid":           "6279acb3696f6e2d31c5fcdd",
					"Name":           "RemoveKubernetesClusterProfileResources",
					"ObjectType":     "workflow.WorkflowInfo",
					"Parent":         nil,
					"ParentTaskInfo": nil,
					"Status":         "RUNNING",
				},
			},
			out: ``,
		},
		{
			msg: "normal workflow update message",
			in: map[string]interface{}{
				"ClassId":         "mo.WebhookResult",
				"ObjectType":      "mo.WebhookResult",
				"EventObjectType": "workflow.WorkflowInfo",
				"Operation":       "Modified",
				"Event": map[string]interface{}{
					"ClassId":        "workflow.WorkflowInfo",
					"Email":          "email@cisco.com",
					"EndTime":        "2022-05-10T00:09:39.362Z",
					"StartTime":      "2022-05-10T00:07:15.844Z",
					"Internal":       false,
					"Moid":           "6279acb3696f6e2d31c5fcdd",
					"Name":           "RemoveKubernetesClusterProfileResources",
					"ObjectType":     "workflow.WorkflowInfo",
					"Parent":         nil,
					"ParentTaskInfo": nil,
					"Status":         "BLAH",
				},
			},
			out: `
## Intersight Workflow Updated

**Name:** [RemoveKubernetesClusterProfileResources](https://www.intersight.com/an/workflow/workflow-infos/6279acb3696f6e2d31c5fcdd)

**Email:** email@cisco.com

**Status:** BLAH

**Start Time:** 2022-05-10T00:07:15.844Z`,
		},
		{
			msg: "finished workflow update message",
			in: map[string]interface{}{
				"ClassId":         "mo.WebhookResult",
				"ObjectType":      "mo.WebhookResult",
				"EventObjectType": "workflow.WorkflowInfo",
				"Operation":       "Modified",
				"Event": map[string]interface{}{
					"ClassId":        "workflow.WorkflowInfo",
					"Email":          "email@cisco.com",
					"EndTime":        "2022-05-10T00:09:39.362Z",
					"StartTime":      "2022-05-10T00:07:15.844Z",
					"Internal":       false,
					"Moid":           "6279acb3696f6e2d31c5fcdd",
					"Name":           "RemoveKubernetesClusterProfileResources",
					"ObjectType":     "workflow.WorkflowInfo",
					"Parent":         nil,
					"ParentTaskInfo": nil,
					"Status":         "COMPLETED",
				},
			},
			out: `
## Intersight Workflow Completed

**Name:** [RemoveKubernetesClusterProfileResources](https://www.intersight.com/an/workflow/workflow-infos/6279acb3696f6e2d31c5fcdd)

**Email:** email@cisco.com

**Status:** COMPLETED

**Start Time:** 2022-05-10T00:07:15.844Z`,
		},
		{
			msg: "internal workflows that end in error should display message",
			in: map[string]interface{}{
				"ClassId":         "mo.WebhookResult",
				"ObjectType":      "mo.WebhookResult",
				"EventObjectType": "workflow.WorkflowInfo",
				"Operation":       "Modified",
				"Event": map[string]interface{}{
					"ClassId":        "workflow.WorkflowInfo",
					"Email":          "email@cisco.com",
					"EndTime":        "2022-05-10T00:09:39.362Z",
					"StartTime":      "2022-05-10T00:07:15.844Z",
					"Internal":       true,
					"Moid":           "6279acb3696f6e2d31c5fcdd",
					"Name":           "RemoveKubernetesClusterProfileResources",
					"ObjectType":     "workflow.WorkflowInfo",
					"Parent":         nil,
					"ParentTaskInfo": nil,
					"Status":         "TERMINATED",
				},
			},
			out: `
## Intersight Workflow Terminated

**Name:** [RemoveKubernetesClusterProfileResources](https://www.intersight.com/an/workflow/workflow-infos/6279acb3696f6e2d31c5fcdd)

**Email:** email@cisco.com

**Status:** TERMINATED

**Start Time:** 2022-05-10T00:07:15.844Z`,
		},
		{
			msg: "unsupported event object types should send a generic message with raw event",
			in: map[string]interface{}{
				"ClassId":         "mo.WebhookResult",
				"ObjectType":      "mo.WebhookResult",
				"EventObjectType": "tam.AdvisoryInfo",
				"Operation":       "Modified",
				"Event": map[string]interface{}{
					"ClassId":        "workflow.WorkflowInfo",
					"Email":          "email@cisco.com",
					"EndTime":        "2022-05-10T00:09:39.362Z",
					"StartTime":      "2022-05-10T00:07:15.844Z",
					"Internal":       true,
					"Moid":           "6279acb3696f6e2d31c5fcdd",
					"Name":           "RemoveKubernetesClusterProfileResources",
					"ObjectType":     "workflow.WorkflowInfo",
					"Parent":         nil,
					"ParentTaskInfo": nil,
					"Status":         "TERMINATED",
				},
			},
			out: `
## Intersight tam.AdvisoryInfo Modified

An Intersight event was received with an event type that I don't support yet, but here is the raw event:

` + "```" + `
{
  "ClassId": "workflow.WorkflowInfo",
  "Email": "email@cisco.com",
  "EndTime": "2022-05-10T00:09:39.362Z",
  "Internal": true,
  "Moid": "6279acb3696f6e2d31c5fcdd",
  "Name": "RemoveKubernetesClusterProfileResources",
  "ObjectType": "workflow.WorkflowInfo",
  "Parent": null,
  "ParentTaskInfo": null,
  "StartTime": "2022-05-10T00:07:15.844Z",
  "Status": "TERMINATED"
}
` + "```",
		},
		{
			msg: "PING events should send an informational message",
			in: map[string]interface{}{
				"Subscription": map[string]interface{}{
					"link":       "https://www.intersight.com/api/v1/notification/AccountSubscriptions/6279fd257375732d3044a48b",
					"ClassId":    "mo.MoRef",
					"Moid":       "6279fd257375732d3044a48b",
					"ObjectType": "notification.AccountSubscription",
				},
				"ObjectType":      "mo.WebhookResult",
				"Event":           nil,
				"AccountMoid":     "59c84e4a16267c0001c23428",
				"DomainGroupMoid": "5b25418d7a7662743465cf72",
				"ClassId":         "mo.WebhookResult",
				"EventObjectType": "",
				"Operation":       "None",
			},
			out: "Intersight webhooks are now being received for subscription [6279fd257375732d3044a48b](https://www.intersight.com/an/settings/webhooks/6279fd257375732d3044a48b/edit/)",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.out, webhookToMessage(test.in), test.msg)
	}
}
