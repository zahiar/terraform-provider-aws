package organizations_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func testAccOrganizationalUnitChildAccountsDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_organizations_organizational_unit_child_accounts.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckOrganizationManagementAccount(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationalUnitChildAccountsDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrGreaterThanValue(dataSourceName, "accounts.#", "0"),
				),
			},
		},
	})
}

const testAccOrganizationalUnitChildAccountsDataSourceConfig_basic = `
data "aws_organizations_organization" "current" {}

data "aws_organizations_organizational_unit_child_accounts" "test" {
  parent_id = data.aws_organizations_organization.current.roots[0].id
}
`
