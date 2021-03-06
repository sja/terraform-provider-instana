package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

//ResourceInstanaApplicationConfig the name of the terraform-provider-instana resource to manage application config
const ResourceInstanaApplicationConfig = "instana_application_config"

const (
	//ApplicationConfigFieldLabel const for the label field of the application config
	ApplicationConfigFieldLabel = "label"
	//ApplicationConfigFieldFullLabel const for the full label field of the application config. The field is computed and contains the label which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level
	ApplicationConfigFieldFullLabel = "full_label"
	//ApplicationConfigFieldScope const for the scope field of the application config
	ApplicationConfigFieldScope = "scope"
	//ApplicationConfigFieldBoundaryScope const for the boundary_scope field of the application config
	ApplicationConfigFieldBoundaryScope = "boundary_scope"
	//ApplicationConfigFieldMatchSpecification const for the match_specification field of the application config
	ApplicationConfigFieldMatchSpecification = "match_specification"
)

var (
	//ApplicationConfigLabel schema for the application config field label
	ApplicationConfigLabel = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The label of the application config",
	}
	//ApplicationConfigFullLabel schema for the application config field full_label
	ApplicationConfigFullLabel = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The the full label field of the application config. The field is computed and contains the label which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	}
	//ApplicationConfigScope schema for the application config field scope
	ApplicationConfigScope = &schema.Schema{
		Type:         schema.TypeString,
		Required:     false,
		Optional:     true,
		Default:      string(restapi.ApplicationConfigScopeIncludeNoDownstream),
		ValidateFunc: validation.StringInSlice(restapi.SupportedApplicationConfigScopes.ToStringSlice(), false),
		Description:  "The scope of the application config",
	}
	//ApplicationConfigBoundaryScope schema for the application config field boundary_scope
	ApplicationConfigBoundaryScope = &schema.Schema{
		Type:         schema.TypeString,
		Required:     false,
		Optional:     true,
		Default:      string(restapi.BoundaryScopeDefault),
		ValidateFunc: validation.StringInSlice(restapi.SupportedBoundaryScopes.ToStringSlice(), false),
		Description:  "The boundary scope of the application config",
	}
	//ApplicationConfigMatchSpecification schema for the application config field match_specification
	ApplicationConfigMatchSpecification = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The match specification of the application config",
	}
)

//NewApplicationConfigResourceHandle creates a new instance of the ResourceHandle for application configs
func NewApplicationConfigResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaApplicationConfig,
		Schema: map[string]*schema.Schema{
			ApplicationConfigFieldLabel:              ApplicationConfigLabel,
			ApplicationConfigFieldFullLabel:          ApplicationConfigFullLabel,
			ApplicationConfigFieldScope:              ApplicationConfigScope,
			ApplicationConfigFieldBoundaryScope:      ApplicationConfigBoundaryScope,
			ApplicationConfigFieldMatchSpecification: ApplicationConfigMatchSpecification,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    applicationConfigSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: applicationConfigStateUpgradeV0,
				Version: 0,
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.ApplicationConfigs() },
		UpdateState:          updateStateForApplicationConfig,
		MapStateToDataObject: mapStateToDataObjectForApplicationConfig,
	}
}

func updateStateForApplicationConfig(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	applicationConfig := obj.(restapi.ApplicationConfig)
	normalizedExpressionString, err := mapAPIModelToNormalizedStringRepresentation(applicationConfig.MatchSpecification.(restapi.MatchExpression))
	if err != nil {
		return err
	}

	d.Set(ApplicationConfigFieldFullLabel, applicationConfig.Label)
	d.Set(ApplicationConfigFieldScope, string(applicationConfig.Scope))
	d.Set(ApplicationConfigFieldBoundaryScope, string(applicationConfig.BoundaryScope))
	d.Set(ApplicationConfigFieldMatchSpecification, normalizedExpressionString)

	d.SetId(applicationConfig.ID)
	return nil
}

func mapAPIModelToNormalizedStringRepresentation(input restapi.MatchExpression) (string, error) {
	mapper := filterexpression.NewMapper()
	expr, err := mapper.FromAPIModel(input)
	if err != nil {
		return "", err
	}
	return expr.Render(), nil
}

func mapStateToDataObjectForApplicationConfig(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	matchSpecification, err := mapExpressionStringToAPIModel(d.Get(ApplicationConfigFieldMatchSpecification).(string))
	if err != nil {
		return restapi.ApplicationConfig{}, err
	}

	label := computeFullApplicationConfigLabelString(d, formatter)
	return restapi.ApplicationConfig{
		ID:                 d.Id(),
		Label:              label,
		Scope:              restapi.ApplicationConfigScope(d.Get(ApplicationConfigFieldScope).(string)),
		BoundaryScope:      restapi.BoundaryScope(d.Get(ApplicationConfigFieldBoundaryScope).(string)),
		MatchSpecification: matchSpecification,
	}, nil
}

func mapExpressionStringToAPIModel(input string) (restapi.MatchExpression, error) {
	parser := filterexpression.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := filterexpression.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func computeFullApplicationConfigLabelString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(ApplicationConfigFieldLabel) {
		return formatter.Format(d.Get(ApplicationConfigFieldLabel).(string))
	}
	return d.Get(ApplicationConfigFieldFullLabel).(string)
}

func applicationConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			ApplicationConfigFieldLabel:              ApplicationConfigLabel,
			ApplicationConfigFieldScope:              ApplicationConfigScope,
			ApplicationConfigFieldMatchSpecification: ApplicationConfigMatchSpecification,
		},
	}
}

func applicationConfigStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState[ApplicationConfigFieldFullLabel] = rawState[ApplicationConfigFieldLabel]
	return rawState, nil
}
