package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/greengrass"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSGreengrassResourceDefinition_basic(t *testing.T) {
	rString := acctest.RandString(8)
	resourceName := "aws_greengrass_resource_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGreengrassResourceDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_basic(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("resource_definition_%s", rString)),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr("aws_greengrass_resource_definition.test", "tags.tagKey", "tagValue"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSGreengrassResourceDefinition_versionNoop(t *testing.T) {
	rString := acctest.RandString(8)
	resourceName := "aws_greengrass_resource_definition.test"

	var resource_version_arn_a string
	var resource_version_arn_b string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGreengrassResourceDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_LocalDevice(rString),
				Check:  testAccAWSGreengrassResourceDefinitionCheckAndGetResourceVersion(resourceName, &resource_version_arn_a),
			},
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_LocalDevice(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSGreengrassResourceDefinitionCheckAndGetResourceVersion(resourceName, &resource_version_arn_b),
					func() resource.TestCheckFunc {
						return func(s *terraform.State) error {
							if resource_version_arn_a != resource_version_arn_b {
								return fmt.Errorf("Resource version ARN %s has changed to %s", resource_version_arn_a, resource_version_arn_b)
							}
							return nil
						}
					}(),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSGreengrassResourceDefinition_versionDiff(t *testing.T) {
	rString := acctest.RandString(8)
	resourceName := "aws_greengrass_resource_definition.test"

	var resource_version_arn_a string
	var resource_version_arn_b string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGreengrassResourceDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_LocalDevice(rString),
				Check:  testAccAWSGreengrassResourceDefinitionCheckAndGetResourceVersion(resourceName, &resource_version_arn_a),
			},
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_LocalVolume(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccAWSGreengrassResourceDefinitionCheckAndGetResourceVersion(resourceName, &resource_version_arn_b),
					func() resource.TestCheckFunc {
						return func(s *terraform.State) error {
							if resource_version_arn_a == resource_version_arn_b {
								return fmt.Errorf("Resource version ARN %s has not changed", resource_version_arn_b)
							}
							return nil
						}
					}(),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSGreengrassResourceDefinition_LocalDevice(t *testing.T) {
	rString := acctest.RandString(8)
	resourceName := "aws_greengrass_resource_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGreengrassResourceDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_LocalDevice(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("resource_definition_%s", rString)),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSGreengrassResourceDefinition_LocalVolume(t *testing.T) {
	rString := acctest.RandString(8)
	resourceName := "aws_greengrass_resource_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGreengrassResourceDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_LocalVolume(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("resource_definition_%s", rString)),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSGreengrassResourceDefinition_S3MachineLearningModel(t *testing.T) {
	rString := acctest.RandString(8)
	resourceName := "aws_greengrass_resource_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGreengrassResourceDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_S3MachineLearningModel(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("resource_definition_%s", rString)),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSGreengrassResourceDefinition_SagemakerMachineLearningModel(t *testing.T) {
	rString := acctest.RandString(8)
	resourceName := "aws_greengrass_resource_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGreengrassResourceDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_SagemakerMachineLearningModel(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("resource_definition_%s", rString)),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSGreengrassResourceDefinition_SecretsManagerSecret(t *testing.T) {
	rString := acctest.RandString(8)
	resourceName := "aws_greengrass_resource_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSGreengrassResourceDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSGreengrassResourceDefinitionConfig_SecretsManagerSecret(rString),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("resource_definition_%s", rString)),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAWSGreengrassResourceDefinitionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).greengrassconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_greengrass_resource_definition" {
			continue
		}

		params := &greengrass.ListResourceDefinitionsInput{
			MaxResults: aws.String("20"),
		}

		out, err := conn.ListResourceDefinitions(params)
		if err != nil {
			return err
		}
		for _, definition := range out.Definitions {
			if *definition.Id == rs.Primary.ID {
				return fmt.Errorf("Expected Greengrass Resource Definition to be destroyed, %s found", rs.Primary.ID)
			}
		}
	}
	return nil
}

func testAccAWSGreengrassResourceDefinitionCheckAndGetResourceVersion(resourceName string, resource_version_arn *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		latest_definition_version_arn, ok := rs.Primary.Attributes["latest_definition_version_arn"]
		if !ok {
			return fmt.Errorf("Latest definition version ARN does not exist for %s", resourceName)
		}

		*resource_version_arn = latest_definition_version_arn

		return nil
	}
}

func testAccAWSGreengrassResourceDefinitionConfig_basic(rString string) string {
	return fmt.Sprintf(`
resource "aws_greengrass_resource_definition" "test" {
  name = "resource_definition_%s"

  tags = {
	"tagKey" = "tagValue"
  }
}
`, rString)
}

func testAccAWSGreengrassResourceDefinitionConfig_LocalDevice(rString string) string {
	return fmt.Sprintf(`
resource "aws_greengrass_resource_definition" "test" {
	name = "resource_definition_%[1]s"
    resource {
        id = "test_id"
        name = "test_name"
        data_container {
            local_device_resource_data {
                source_path = "/dev/source"

                group_owner_setting {
                    auto_add_group_owner = false
                    group_owner = "user"
                }
            }
        }
	}
}
`, rString)
}

func testAccAWSGreengrassResourceDefinitionConfig_LocalVolume(rString string) string {
	return fmt.Sprintf(`
resource "aws_greengrass_resource_definition" "test" {
	name = "resource_definition_%[1]s"
    resource {
        id = "test_id"
        name = "test_name"
        data_container {

            local_volume_resource_data {
                source_path = "/dev/source"
                destination_path = "/destination"

                group_owner_setting {
                    auto_add_group_owner = false
                    group_owner = "user"
                }
            }
        }
	}
}
`, rString)
}

func testAccAWSGreengrassResourceDefinitionConfig_S3MachineLearningModel(rString string) string {
	return fmt.Sprintf(`
resource "aws_greengrass_resource_definition" "test" {
	name = "resource_definition_%[1]s"
    resource {
        id = "test_id"
        name = "test_name"
        data_container {
            s3_machine_learning_model_resource_data {
                s3_uri = "s3://bucket/key.zip"
                destination_path = "/destination"
            }
        }
	}
}
`, rString)
}

func testAccAWSGreengrassResourceDefinitionConfig_SagemakerMachineLearningModel(rString string) string {
	return fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_greengrass_resource_definition" "test" {
	name = "resource_definition_%[1]s"
    resource {
        id = "test_id"
        name = "test_name"
        data_container {
            sagemaker_machine_learning_model_resource_data {
                sagemaker_job_arn = "arn:aws:sagemaker:us-west-2:${data.aws_caller_identity.current.account_id}:training-job/xgboost-2018-06-05-17-19-32-703"
                destination_path = "/destination"
            }
        }
    }
}
`, rString)
}

func testAccAWSGreengrassResourceDefinitionConfig_SecretsManagerSecret(rString string) string {
	return fmt.Sprintf(`
resource "aws_greengrass_resource_definition" "test" {
	name = "resource_definition_%[1]s"
    resource {
        id = "test_id"
        name = "test_name"
        data_container {
            secrets_manager_secret_resource_data {
                secret_arn = "arn:aws:secretsmanager:us-west-2:123456789012:secret:greengrass-TwilioAuthToken-ntSlp6"
                additional_staging_labels_to_download = [
                    "label1",
                    "label2",
                ]
            }
        }
	}
}
`, rString)
}
