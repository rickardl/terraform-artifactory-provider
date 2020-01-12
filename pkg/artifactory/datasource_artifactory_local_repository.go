package artifactory

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rickardl/go-artifactory/v2/artifactory"
)

func dataSourceArtifactoryLocalRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataLocalRepositoryRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"package_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Computed: true,
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
			"debian_trivial_layout": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"checksum_policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_unique_tags": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"snapshot_version_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"suppress_pom_consistency_checks": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"blacked_out": {
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
			"archive_browsing_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"calculate_yum_metadata": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"yum_root_depth": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"docker_api_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_file_lists_indexing": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"xray_index": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataLocalRepositoryRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*artifactory.Artifactory)

	key := d.Get("key").(string)
	log.Printf("[DEBUG] Reading Local Repository with Key: %s", key)
	repo, resp, err := c.V1.Repositories.GetLocal(context.Background(), key)

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
	} else if err == nil {
		hasErr := false
		logError := cascadingErr(&hasErr)
		logError(d.Set("key", repo.Key))
		logError(d.Set("package_type", repo.PackageType))
		logError(d.Set("description", repo.Description))
		logError(d.Set("notes", repo.Notes))
		logError(d.Set("includes_pattern", repo.IncludesPattern))
		logError(d.Set("excludes_pattern", repo.ExcludesPattern))
		logError(d.Set("repo_layout_ref", repo.RepoLayoutRef))
		logError(d.Set("debian_trivial_layout", repo.DebianTrivialLayout))
		logError(d.Set("max_unique_tags", repo.MaxUniqueTags))
		logError(d.Set("blacked_out", repo.BlackedOut))
		logError(d.Set("archive_browsing_enabled", repo.ArchiveBrowsingEnabled))
		logError(d.Set("calculate_yum_metadata", repo.CalculateYumMetadata))
		logError(d.Set("yum_root_depth", repo.YumRootDepth))
		logError(d.Set("docker_api_version", repo.DockerApiVersion))
		logError(d.Set("enable_file_lists_indexing", repo.EnableFileListsIndexing))
		logError(d.Set("property_sets", schema.NewSet(schema.HashString, castToInterfaceArr(*repo.PropertySets))))
		logError(d.Set("handle_releases", repo.HandleReleases))
		logError(d.Set("handle_snapshots", repo.HandleSnapshots))
		logError(d.Set("checksum_policy_type", repo.ChecksumPolicyType))
		logError(d.Set("max_unique_snapshots", repo.MaxUniqueSnapshots))
		logError(d.Set("snapshot_version_behavior", repo.SnapshotVersionBehavior))
		logError(d.Set("suppress_pom_consistency_checks", repo.SuppressPomConsistencyChecks))
		logError(d.Set("xray_index", repo.XrayIndex))

		if hasErr {
			return fmt.Errorf("failed to marshal group")
		}
		d.SetId(*repo.Key)

	}

	return err
}
