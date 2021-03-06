package restapi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/utils"
)

//AlertsResourcePath path to Alerts resource of Instana RESTful API
const AlertsResourcePath = EventSettingsBasePath + "/alerts"

//AlertEventType type definition of EventTypes of an Instana Alert
type AlertEventType string

//Equals checks if the alert event type is equal to the provided alert event type. It compares the string representation of both case insensitive
func (t AlertEventType) Equals(other AlertEventType) bool {
	return strings.ToLower(string(t)) == strings.ToLower(string(other))
}

const (
	//IncidentAlertEventType constant value for alert event type incident
	IncidentAlertEventType = AlertEventType("incident")
	//CriticalAlertEventType constant value for alert event type critical
	CriticalAlertEventType = AlertEventType("critical")
	//WarningAlertEventType constant value for alert event type warning
	WarningAlertEventType = AlertEventType("warning")
	//ChangeAlertEventType constant value for alert event type change
	ChangeAlertEventType = AlertEventType("change")
	//OnlineAlertEventType constant value for alert event type online
	OnlineAlertEventType = AlertEventType("online")
	//OfflineAlertEventType constant value for alert event type offline
	OfflineAlertEventType = AlertEventType("offline")
	//NoneAlertEventType constant value for alert event type none
	NoneAlertEventType = AlertEventType("none")
	//AgentMonitoringIssueEventType constant value for alert event type none
	AgentMonitoringIssueEventType = AlertEventType("agent_monitoring_issue")
)

//SupportedAlertEventTypes list of supported alert event types of Instana API
var SupportedAlertEventTypes = []AlertEventType{
	IncidentAlertEventType,
	CriticalAlertEventType,
	WarningAlertEventType,
	ChangeAlertEventType,
	OnlineAlertEventType,
	OfflineAlertEventType,
	NoneAlertEventType,
	AgentMonitoringIssueEventType,
}

//IsSupportedAlertEventType checks if the given alert type is supported by Instana API
func IsSupportedAlertEventType(t AlertEventType) bool {
	for _, supported := range SupportedAlertEventTypes {
		if supported.Equals(t) {
			return true
		}
	}
	return false
}

//EventFilteringConfiguration type definiton of an EventFilteringConfiguration of a AlertingConfiguration of the Instana ReST AOI
type EventFilteringConfiguration struct {
	Query      *string          `json:"query"`
	RuleIDs    []string         `json:"ruleIds"`
	EventTypes []AlertEventType `json:"eventTypes"`
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (c EventFilteringConfiguration) Validate() error {
	if c.Query != nil && len(*c.Query) > 2048 {
		return errors.New("Query of EventFilterConfig not valid; Maximum length of Query is 2048 characters")
	}

	numberRuleIDs := len(c.RuleIDs)
	numberEventType := len(c.EventTypes)
	if (numberRuleIDs == 0 && numberEventType == 0) || (numberRuleIDs > 0 && numberEventType > 0) {
		return errors.New("Invalid EventFilterConfig; Either RuleIDs or EventTypes must be configured")
	}

	if numberRuleIDs > 1024 {
		return errors.New("Invalid EventFilterConfig; Maximum number of RuleIDs is 1024")
	}

	if !utils.StringSliceElementsAreUnique(c.RuleIDs) {
		return errors.New("Invalid EventFilterConfig; RuleIDs must be unique")
	}

	if numberEventType > len(SupportedAlertEventTypes) {
		return fmt.Errorf("Invalid EventFilterConfig; Maximum number of EventTypes is %d", len(SupportedAlertEventTypes))
	}

	if !utils.StringSliceElementsAreUnique(eventTypeSliceToStringSlice(c.EventTypes)) {
		return errors.New("Invalid EventFilterConfig; EventTypes must be unique")
	}

	for _, t := range c.EventTypes {
		if !IsSupportedAlertEventType(t) {
			return fmt.Errorf("%s is not a supported EventType", t)
		}
	}
	return nil
}

func eventTypeSliceToStringSlice(input []AlertEventType) []string {
	numberElements := len(input)
	if numberElements == 0 {
		return []string{}
	}
	result := make([]string, numberElements)
	for i, v := range input {
		result[i] = string(v)
	}
	return result
}

//AlertingConfiguration type definition of an Alertinng Configruation in Instana REST API
type AlertingConfiguration struct {
	ID                          string                      `json:"id"`
	AlertName                   string                      `json:"alertName"`
	IntegrationIDs              []string                    `json:"integrationIds"`
	EventFilteringConfiguration EventFilteringConfiguration `json:"eventFilteringConfiguration"`
}

//GetID implemention of the interface InstanaDataObject
func (c AlertingConfiguration) GetID() string {
	return c.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (c AlertingConfiguration) Validate() error {
	if utils.IsBlank(c.ID) {
		return errors.New("ID is missing")
	}
	if utils.IsBlank(c.AlertName) {
		return errors.New("AlertName is missing")
	}
	if len(c.AlertName) > 256 {
		return errors.New("AlertName not valid; Maximum length of AlertName is 256 characters")
	}
	if len(c.IntegrationIDs) > 1024 {
		return errors.New("Too many IntegrationIDs; Maximum number of IntegrationIds is 1024")
	}

	if !utils.StringSliceElementsAreUnique(c.IntegrationIDs) {
		return errors.New("IntegrationIDs must be unique")
	}
	return c.EventFilteringConfiguration.Validate()
}
