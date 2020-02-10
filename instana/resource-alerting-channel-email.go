package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//AlertingChannelEmailFieldEmails const for the emails field of the alerting channel
	AlertingChannelEmailFieldEmails = "emails"
	//ResourceInstanaAlertingChannelEmail the name of the terraform-provider-instana resource to manage alerting channels of type email
	ResourceInstanaAlertingChannelEmail = "instana_alerting_channel_email"
)

//NewAlertingChannelEmailResourceHandle creates the resource handle for Alerting Channels of type Email
func NewAlertingChannelEmailResourceHandle() ResourceHandle {
	return &alertingChannelEmailResourceHandle{}
}

type alertingChannelEmailResourceHandle struct {
}

func (h *alertingChannelEmailResourceHandle) GetResourceFrom(api restapi.InstanaAPI) restapi.RestResource {
	return api.AlertingChannels()
}

func (h *alertingChannelEmailResourceHandle) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		AlertingChannelFieldName:     alertingChannelNameSchemaField,
		AlertingChannelFieldFullName: alertingChannelFullNameSchemaField,
		AlertingChannelEmailFieldEmails: &schema.Schema{
			Type:     schema.TypeList,
			MinItems: 1,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required:    true,
			Description: "The list of emails of the Email alerting channel",
		},
	}
}

func (h *alertingChannelEmailResourceHandle) SchemaVersion() int {
	return 0
}

func (h *alertingChannelEmailResourceHandle) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (h *alertingChannelEmailResourceHandle) ResourceName() string {
	return ResourceInstanaAlertingChannelEmail
}

func (h *alertingChannelEmailResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	alertingChannel := obj.(restapi.AlertingChannel)
	emails := alertingChannel.Emails
	d.Set(AlertingChannelFieldFullName, alertingChannel.Name)
	d.Set(AlertingChannelEmailFieldEmails, emails)
	d.SetId(alertingChannel.ID)
	return nil
}

func (h *alertingChannelEmailResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingChannelNameString(d, formatter)
	return restapi.AlertingChannel{
		ID:     d.Id(),
		Name:   name,
		Kind:   restapi.EmailChannelType,
		Emails: ReadStringArrayParameterFromResource(d, AlertingChannelEmailFieldEmails),
	}, nil
}
