package iswbx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsUpdateCleanup(t *testing.T) {
	var wh map[string]interface{}
	var tm time.Time

	wh = map[string]interface{}{
		"ClassId":         "mo.WebhookResult",
		"ObjectType":      "mo.WebhookResult",
		"EventObjectType": "workflow.WorkflowInfo",
		"Operation":       "Modified",
		"Event": map[string]interface{}{
			"ClassId":        "workflow.WorkflowInfo",
			"Email":          "email@cisco.com",
			"EndTime":        "2022-05-10T00:09:39.362Z",
			"StartTime":      "2022-05-10T00:07:15.844Z",
			"CleanupTime":    "2022-05-13T04:14:01.059Z",
			"Internal":       true,
			"Moid":           "6279acb3696f6e2d31c5fcdd",
			"Name":           "RemoveKubernetesClusterProfileResources",
			"ObjectType":     "workflow.WorkflowInfo",
			"Parent":         nil,
			"ParentTaskInfo": nil,
			"Status":         "TERMINATED",
		},
	}
	tm = time.Date(2022, 5, 13, 4, 15, 36, 0, time.UTC)
	assert.True(t, isWorkflowUpdateCleanup(wh, tm))

	tm = time.Date(2022, 5, 13, 4, 0, 36, 0, time.UTC)
	assert.False(t, isWorkflowUpdateCleanup(wh, tm))
}
