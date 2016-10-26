package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestRegisterAccountLoginWidget(c *C) {
	widget := scopes.NewPreviewWidget("buttons", "actions")
	widget.AddAttributeValue("actions", map[string]interface{}{
		"id":    "install-snap",
		"label": "Install",
	})
	scopes.RegisterAccountLoginWidget(&widget,
		"ubuntuone", "ubuntuone", "ubuntuone",
		scopes.PostLoginContinueActivation,
		scopes.PostLoginDoNothing)

	result, ok := widget["online_account_details"]
	c.Check(ok, Equals, true)
	c.Check(result, Equals, map[string]interface{}{
		"scope_id":            "",
		"service_name":        "ubuntuone",
		"service_type":        "ubuntuone",
		"provider_name":       "ubuntuone",
		"login_passed_action": scopes.PostLoginContinueActivation,
		"login_failed_action": scopes.PostLoginDoNothing,
	})
}
