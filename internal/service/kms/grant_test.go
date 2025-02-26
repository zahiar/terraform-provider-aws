package kms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfkms "github.com/hashicorp/terraform-provider-aws/internal/service/kms"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccKMSGrant_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_kms_grant.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_basic(rName, "\"Encrypt\", \"Decrypt\""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "operations.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "operations.*", "Encrypt"),
					resource.TestCheckTypeSetElemAttr(resourceName, "operations.*", "Decrypt"),
					resource.TestCheckResourceAttrPair(resourceName, "grantee_principal", "aws_iam_role.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "key_id", "aws_kms_key.test", "key_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"grant_token", "retire_on_delete"},
			},
		},
	})
}

func TestAccKMSGrant_withConstraints(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_kms_grant.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_constraints(rName, "encryption_context_equals", `foo = "bar"
                        baz = "kaz"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "constraints.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "constraints.*", map[string]string{
						"encryption_context_equals.%":   "2",
						"encryption_context_equals.baz": "kaz",
						"encryption_context_equals.foo": "bar",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"grant_token", "retire_on_delete"},
			},
			{
				Config: testAccGrantConfig_constraints(rName, "encryption_context_subset", `foo = "bar"
			            baz = "kaz"`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "constraints.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "constraints.*", map[string]string{
						"encryption_context_subset.%":   "2",
						"encryption_context_subset.baz": "kaz",
						"encryption_context_subset.foo": "bar",
					}),
				),
			},
		},
	})
}

func TestAccKMSGrant_withRetiringPrincipal(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_kms_grant.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_retiringPrincipal(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "retiring_principal", "aws_iam_role.test", "arn"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"grant_token", "retire_on_delete"},
			},
		},
	})
}

func TestAccKMSGrant_bare(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_kms_grant.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_bare(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					resource.TestCheckNoResourceAttr(resourceName, "name"),
					resource.TestCheckNoResourceAttr(resourceName, "constraints.#"),
					resource.TestCheckNoResourceAttr(resourceName, "retiring_principal"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"grant_token", "retire_on_delete"},
			},
		},
	})
}

func TestAccKMSGrant_arn(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_kms_grant.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_arn(rName, "\"Encrypt\", \"Decrypt\""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "operations.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "operations.*", "Encrypt"),
					resource.TestCheckTypeSetElemAttr(resourceName, "operations.*", "Decrypt"),
					resource.TestCheckResourceAttrPair(resourceName, "grantee_principal", "aws_iam_role.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "key_id", "aws_kms_key.test", "arn"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"grant_token", "retire_on_delete"},
			},
		},
	})
}

func TestAccKMSGrant_asymmetricKey(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_kms_grant.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_asymmetricKey(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"grant_token", "retire_on_delete"},
			},
		},
	})
}

func TestAccKMSGrant_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_kms_grant.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_basic(rName, "\"Encrypt\", \"Decrypt\""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfkms.ResourceGrant(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKMSGrant_crossAccountARN(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_kms_grant.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckAlternateAccount(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
		CheckDestroy:             testAccCheckGrantDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGrantConfig_crossAccountARN(rName, "\"Encrypt\", \"Decrypt\""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrantExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "operations.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "operations.*", "Encrypt"),
					resource.TestCheckTypeSetElemAttr(resourceName, "operations.*", "Decrypt"),
					resource.TestCheckResourceAttrPair(resourceName, "grantee_principal", "aws_iam_role.test", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "key_id", "aws_kms_key.test", "arn"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"grant_token", "retire_on_delete"},
			},
		},
	})
}

func testAccCheckGrantDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).KMSConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_kms_grant" {
				continue
			}

			keyID, grantID, err := tfkms.GrantParseResourceID(rs.Primary.ID)

			if err != nil {
				return err
			}

			_, err = tfkms.FindGrantByTwoPartKey(ctx, conn, keyID, grantID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("KMS Grant still exists: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckGrantExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No KMS Grant ID is set")
		}

		keyID, grantID, err := tfkms.GrantParseResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).KMSConn()

		_, err = tfkms.FindGrantByTwoPartKey(ctx, conn, keyID, grantID)

		return err
	}
}

func testAccGrantConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

data "aws_iam_policy_document" "test" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "test" {
  name               = %[1]q
  path               = "/service-role/"
  assume_role_policy = data.aws_iam_policy_document.test.json
}
`, rName)
}

func testAccGrantConfig_basic(rName string, operations string) string {
	return acctest.ConfigCompose(testAccGrantConfig_base(rName), fmt.Sprintf(`
resource "aws_kms_grant" "test" {
  name              = %[1]q
  key_id            = aws_kms_key.test.key_id
  grantee_principal = aws_iam_role.test.arn
  operations        = [%[2]s]
}
`, rName, operations))
}

func testAccGrantConfig_constraints(rName string, constraintName string, encryptionContext string) string {
	return acctest.ConfigCompose(testAccGrantConfig_base(rName), fmt.Sprintf(`
resource "aws_kms_grant" "test" {
  name              = %[1]q
  key_id            = aws_kms_key.test.key_id
  grantee_principal = aws_iam_role.test.arn
  operations        = ["RetireGrant", "DescribeKey"]

  constraints {
    %[2]s = {
      %[3]s
    }
  }
}
`, rName, constraintName, encryptionContext))
}

func testAccGrantConfig_retiringPrincipal(rName string) string {
	return acctest.ConfigCompose(testAccGrantConfig_base(rName), fmt.Sprintf(`
resource "aws_kms_grant" "test" {
  name               = %[1]q
  key_id             = aws_kms_key.test.key_id
  grantee_principal  = aws_iam_role.test.arn
  operations         = ["ReEncryptTo", "CreateGrant"]
  retiring_principal = aws_iam_role.test.arn
}
`, rName))
}

func testAccGrantConfig_bare(rName string) string {
	return acctest.ConfigCompose(testAccGrantConfig_base(rName), `
resource "aws_kms_grant" "test" {
  key_id            = aws_kms_key.test.key_id
  grantee_principal = aws_iam_role.test.arn
  operations        = ["ReEncryptTo", "CreateGrant"]
}
`)
}

func testAccGrantConfig_arn(rName string, operations string) string {
	return acctest.ConfigCompose(testAccGrantConfig_base(rName), fmt.Sprintf(`
resource "aws_kms_grant" "test" {
  name              = %[1]q
  key_id            = aws_kms_key.test.arn
  grantee_principal = aws_iam_role.test.arn
  operations        = [%[2]s]
}
`, rName, operations))
}

func testAccGrantConfig_asymmetricKey(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_grant" "test" {
  name              = %[1]q
  key_id            = aws_kms_key.test.key_id
  grantee_principal = aws_iam_role.test.arn
  operations        = ["GetPublicKey", "Sign", "Verify"]
}

resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7

  key_usage                = "SIGN_VERIFY"
  customer_master_key_spec = "RSA_2048"
}

data "aws_iam_policy_document" "test" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "test" {
  name               = %[1]q
  path               = "/service-role/"
  assume_role_policy = data.aws_iam_policy_document.test.json
}
`, rName)
}

func testAccGrantConfig_crossAccountARN(rName string, operations string) string {
	return acctest.ConfigCompose(acctest.ConfigAlternateAccountProvider(), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

data "aws_iam_policy_document" "test" {
  provider = "awsalternate"

  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "test" {
  provider = "awsalternate"

  name               = %[1]q
  path               = "/service-role/"
  assume_role_policy = data.aws_iam_policy_document.test.json
}

resource "aws_kms_grant" "test" {
  name              = %[1]q
  key_id            = aws_kms_key.test.arn
  grantee_principal = aws_iam_role.test.arn
  operations        = [%[2]s]
}
`, rName, operations))
}
