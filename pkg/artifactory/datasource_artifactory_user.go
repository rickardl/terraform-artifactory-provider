package artifactory

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rickardl/go-artifactory/v2/artifactory"
)

func dataSourceArtifactoryUser() *schema.Resource {
	return &schema.Resource{
		Read: dataUserRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"profile_updatable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disable_ui_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"internal_password_disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"groups": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				StateFunc: func(value interface{}) string {
					// Avoid storing the actual value in the state and instead store the hash of it
					return hashString(value.(string))
				},
			},
		},
	}
}

func dataUserRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Artifactory)

	user, resp, err := c.V1.Security.GetUser(context.Background(), d.Id())
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}
	d.SetId(*user.Name)

	return packUser(user, d)
}
