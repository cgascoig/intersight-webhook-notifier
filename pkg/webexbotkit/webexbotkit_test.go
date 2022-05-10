package webexbotkit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageTrim(t *testing.T) {
	out := trimMessageText("Intersight subscription status", "Intersight Notifier-dev")
	assert.Equal(t, "subscription status", out)

	out = trimMessageText("Intersight Notifier-dev subscription status", "Intersight Notifier-dev")
	assert.Equal(t, "subscription status", out)

	out = trimMessageText(" subscription status", "Intersight Notifier-dev")
	assert.Equal(t, "subscription status", out)
}
