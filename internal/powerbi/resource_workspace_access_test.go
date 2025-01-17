package powerbi

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/codecutout/terraform-provider-powerbi/internal/powerbiapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccWorkspaceAccess_basic(t *testing.T) {
	workspaceSuffix := acctest.RandString(6)
	secondaryUsername := os.Getenv("POWERBI_SECONDARY_USERNAME")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if secondaryUsername == "" {
				t.Fatal("POWERBI_SECONDARY_USERNAME must be set for workspace access acceptance tests")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the resource
			{
				Config: fmt.Sprintf(`
				resource "powerbi_workspace" "test" {
					name = "Acceptance Test Workspace %s"
				}

				resource "powerbi_workspace_access" "test" {
					workspace_id = "${powerbi_workspace.test.id}"
					group_user_access_right = "Admin"
					email_address = "%s"
					principal_type = "User"
				}
				`, workspaceSuffix, secondaryUsername),
				Check: resource.ComposeTestCheckFunc(
					testCheckGroupUserExistsInWorkspace("powerbi_workspace.test", secondaryUsername),
					resource.TestCheckResourceAttrSet("powerbi_workspace_access.test", "id"),
					resource.TestCheckResourceAttrSet("powerbi_workspace_access.test", "workspace_id"),
					resource.TestCheckResourceAttr("powerbi_workspace_access.test", "id", fmt.Sprintf("Acceptance Test Workspace %s/%s", workspaceSuffix, secondaryUsername)),
				),
			},
			// final step checks importing the current state we reached in the step above
			{
				ResourceName:      "powerbi_workspace_access.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkspaceAccess_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerbi_workspace_access" "test" {
					workspace_id = "validation-should-fail-before-using-this"
					group_user_access_right = "Admin"
					email_address = "validation-should-fail-before-using-this"
					principal_type = "User"
				}
				`,
				ExpectError: regexp.MustCompile("config is invalid:.*email_address.*"),
			},
			{
				Config: `
				resource "powerbi_workspace_access" "test" {
					workspace_id = "validation-should-fail-before-using-this"
					group_user_access_right = "Admin"
					email_address = "user@mailserver"
				}
				`,
				ExpectError: regexp.MustCompile("config is invalid:.*principal_type.*"),
			},
			{
				Config: `
				resource "powerbi_workspace_access" "test" {
					workspace_id = "validation-should-fail-before-using-this"
					group_user_access_right = "Admin"
					identifier = "validation-should-fail-before-using-this"
					principal_type = "not-valid-type"
				}
				`,
				ExpectError: regexp.MustCompile("config is invalid:.*principal_type.*"),
			},
			{
				Config: `
				resource "powerbi_workspace_access" "test" {
					workspace_id = "validation-should-fail-before-using-this"
					group_user_access_right = "not-valid-access-right"
					identifier = "validation-should-fail-before-using-this"
					principal_type = "App"
				}
				`,
				ExpectError: regexp.MustCompile("config is invalid:.*group_user_access_right.*"),
			},
		},
	})
}

func TestAccWorkspaceAccess_skew(t *testing.T) {
	var workspaceUserID string
	var groupID string
	workspaceSuffix := acctest.RandString(6)
	secondaryUsername := os.Getenv("POWERBI_SECONDARY_USERNAME")

	config := fmt.Sprintf(`
	resource "powerbi_workspace" "test" {
		name = "Acceptance Test Workspace %s"
	}

	resource "powerbi_workspace_access" "test" {
		workspace_id = "${powerbi_workspace.test.id}"
		group_user_access_right = "Admin"
		email_address = "%s"
		principal_type = "User"
	}
	`, workspaceSuffix, secondaryUsername)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			if secondaryUsername == "" {
				t.Fatal("POWERBI_SECONDARY_USERNAME must be set for workspace access acceptance tests")
			}
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPowerbiWorkspaceDestroy,
		Steps: []resource.TestStep{
			// first step creates the resource
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					set("powerbi_workspace_access.test", "id", &workspaceUserID),
					set("powerbi_workspace.test", "id", &groupID),
				),
			},
			// second step skew new access right
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.UpdateGroupUser(groupID, powerbiapi.UpdateGroupUserRequest{
						//Identifier:           secondaryUsername,
						//PrincipalType:        "User",
						EmailAddress:         "%s",
						GroupUserAccessRight: "Member",
					})
					client.RefreshUserPermissions()
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckGroupUserExistsInWorkspace("powerbi_workspace.test", secondaryUsername),
					resource.TestCheckResourceAttr("powerbi_workspace_access.test", "group_user_access_right", "Admin"),
				),
			},
			// third step skew by deleting user
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*powerbiapi.Client)
					client.DeleteUserInGroup(groupID, workspaceUserID)
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckGroupUserExistsInWorkspace("powerbi_workspace.test", secondaryUsername),
				),
			},
		},
	})
}

func testCheckGroupUserExistsInWorkspace(workspaceResourceName string, expectedIdentifier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupID, err := getResourceID(s, workspaceResourceName)
		if err != nil {
			return err
		}

		var userObjFound bool

		client := testAccProvider.Meta().(*powerbiapi.Client)
		groupUsers, err := client.GetGroupUsers(groupID)
		if err != nil {
			return err
		}

		if len(groupUsers.Value) >= 1 {
			for _, userObj := range groupUsers.Value {
				if userObj.Identifier == expectedIdentifier {
					userObjFound = true
				}
			}
		}
		if userObjFound != true {
			return fmt.Errorf("Expecting groupusers %v in workspace %v. Not found", expectedIdentifier, groupID)
		}
		return nil
	}
}
