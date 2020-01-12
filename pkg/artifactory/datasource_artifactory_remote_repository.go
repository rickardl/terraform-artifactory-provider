package artifactory

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rickardl/go-artifactory/v2/artifactory"
)

func dataSourceArtifactoryRemoteRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataRemoteRepositoryRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"package_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return old == fmt.Sprintf("%s (local file cache)", new)
				},
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"includes_pattern": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"excludes_pattern": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"repo_layout_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"handle_releases": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"handle_snapshots": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"max_unique_snapshots": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"suppress_pom_consistency_checks": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				StateFunc: getMD5Hash,
			},
			"proxy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"remote_repo_checksum_policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"generate-if-absent",
					"fail",
					"ignore-and-generate",
					"pass-thru",
				}, false),
			},
			"hard_fail": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"offline": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"blacked_out": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"store_artifacts_locally": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"socket_timeout_millis": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"local_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retrieval_cache_period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"missed_cache_period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"unused_artifacts_cleanup_period_hours": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"fetch_jars_eagerly": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"fetch_sources_eagerly": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"share_configuration": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"synchronize_properties": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"block_mismatching_mime_types": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"property_sets": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"allow_any_host_auth": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_cookie_management": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"client_tls_certificate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pypi_registry_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bower_registry_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bypass_head_requests": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_token_authentication": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"xray_index": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"vcs_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vcs_git_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vcs_git_download_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"feed_context_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"nuget"},
			},
			"download_context_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"nuget"},
			},
			"v3_feed_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"nuget"},
			},
			"nuget": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				MinItems:      1,
				Deprecated:    "Since Artifactory 6.9.0+ (provider 1.6). Use /api/v2 endpoint",
				ConflictsWith: []string{"feed_context_path", "download_context_path", "v3_feed_url"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"feed_context_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"download_context_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"v3_feed_url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataRemoteRepositoryRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Artifactory)

	key := d.Get("key").(string)
	log.Printf("[DEBUG] Reading Local Repository with Key: %s", key)

	repo, resp, err := c.V1.Repositories.GetRemote(context.Background(), key)
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}
	d.SetId(*repo.Key)
	return packRemoteRepo(repo, d)
}
