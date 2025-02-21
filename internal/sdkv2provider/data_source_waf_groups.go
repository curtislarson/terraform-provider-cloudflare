package sdkv2provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	cloudflare "github.com/curtislarson/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceCloudflareWAFGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCloudflareWAFGroupsRead,

		Schema: map[string]*schema.Schema{
			consts.ZoneIDSchemaKey: {
				Description: "The zone identifier to target for the resource.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"package_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
						},
					},
				},
			},

			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rules_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"modified_rules_count": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"package_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudflareWAFGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*cloudflare.API)
	zoneID := d.Get(consts.ZoneIDSchemaKey).(string)

	// Prepare the filters to be applied to the search
	filter, err := expandFilterWAFGroups(d.Get("filter"))
	if err != nil {
		return diag.FromErr(err)
	}

	// If no package ID is given, we will consider all for the zone
	packageID := d.Get("package_id").(string)
	var pkgList []cloudflare.WAFPackage
	if packageID == "" {
		var err error
		tflog.Debug(ctx, fmt.Sprintf("Reading WAF Packages"))
		pkgList, err = client.ListWAFPackages(ctx, zoneID)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		pkgList = append(pkgList, cloudflare.WAFPackage{ID: packageID})
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading WAF Groups"))
	groupIds := make([]string, 0)
	groupDetails := make([]interface{}, 0)
	for _, pkg := range pkgList {
		groupList, err := client.ListWAFGroups(ctx, zoneID, pkg.ID)
		if err != nil {
			return diag.FromErr(err)
		}

		for _, group := range groupList {
			if filter.Name != nil && !filter.Name.Match([]byte(group.Name)) {
				continue
			}

			if filter.Mode != "" && filter.Mode != group.Mode {
				continue
			}

			groupDetails = append(groupDetails, map[string]interface{}{
				"id":                   group.ID,
				"name":                 group.Name,
				"description":          group.Description,
				"mode":                 group.Mode,
				"rules_count":          group.RulesCount,
				"modified_rules_count": group.ModifiedRulesCount,
				"package_id":           pkg.ID,
			})
			groupIds = append(groupIds, group.ID)
		}
	}

	err = d.Set("groups", groupDetails)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting WAF groups: %w", err))
	}

	d.SetId(stringListChecksum(groupIds))
	return nil
}

func expandFilterWAFGroups(d interface{}) (*searchFilterWAFGroups, error) {
	cfg := d.([]interface{})
	filter := &searchFilterWAFGroups{}
	if len(cfg) == 0 || cfg[0] == nil {
		return filter, nil
	}

	m := cfg[0].(map[string]interface{})
	name, ok := m["name"]
	if ok {
		match, err := regexp.Compile(name.(string))
		if err != nil {
			return nil, err
		}

		filter.Name = match
	}

	mode, ok := m["mode"]
	if ok {
		filter.Mode = mode.(string)
	}

	return filter, nil
}

type searchFilterWAFGroups struct {
	Name *regexp.Regexp
	Mode string
}
