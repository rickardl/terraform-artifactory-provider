package artifactory

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
)

func testAccDataSourceLocalRepoConfig_basic(randInt int) string {
	return fmt.Sprintf(`
	key 	     = "terraform-local-test-repo-basic-%d"
	package_type = "docker"
`, randInt)
}

func TestAccDataLocalRepository_basic(t *testing.T) {

	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLocalRepoConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("artifactory_local_repository.terraform-local-test-repo-basic", "key", "terraform-local-test-repo-basic"),
					resource.TestCheckResourceAttr("artifactory_local_repository.terraform-local-test-repo-basic", "package_type", "docker"),
				),
			},
		},
	})
}
