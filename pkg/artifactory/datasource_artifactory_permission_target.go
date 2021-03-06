package artifactory

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rickardl/go-artifactory/v2/artifactory"
	v2 "github.com/rickardl/go-artifactory/v2/artifactory/v2"
)

func dataSourceArtifactoryPermissionTarget() *schema.Resource {

	actionSchema := schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Set:      hashPrincipal,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"permissions": {
					Type: schema.TypeSet,
					Elem: &schema.Schema{
						Type: schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							v2.PERM_READ,
							v2.PERM_ANNOTATE,
							v2.PERM_WRITE,
							v2.PERM_DELETE,
							v2.PERM_MANAGE,
						}, false),
					},
					Set:      schema.HashString,
					Required: true,
				},
			},
		},
	}

	principalSchema := schema.Schema{
		Type:          schema.TypeList,
		ConflictsWith: []string{"repositories"},
		Optional:      true,
		MaxItems:      1,
		MinItems:      1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"includes_pattern": {
					Type:     schema.TypeSet,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Set:      schema.HashString,
					Optional: true,
				},
				"excludes_pattern": {
					Type:     schema.TypeSet,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Set:      schema.HashString,
					Optional: true,
				},
				"repositories": {
					Type:     schema.TypeSet,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Set:      schema.HashString,
					Required: true,
				},
				"actions": {
					Type:     schema.TypeList,
					MaxItems: 1,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"users":  &actionSchema,
							"groups": &actionSchema,
						},
					},
					Optional: true,
				},
			},
		},
	}

	return &schema.Resource{
		Read: dataPermissionTargetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repo":  &principalSchema,
			"build": &principalSchema,
		},
	}
}

func dataPermissionTargetRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Artifactory)

	name := d.Get("name").(string)
	log.Printf("[DEBUG] Reading Perssmion Target with name: %s", name)

	permissionTarget, resp, err := c.V2.Security.GetPermissionTarget(context.Background(), name)
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	d.SetId(*permissionTarget.Name)
	return packPermissionTarget(permissionTarget, d)
}
