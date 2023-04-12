/*
Copyright 2023 patrick hermann.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
)

func ConvertVaultSecretsInParameters(parameters []string) (updatedParameters []string) {

	for _, parameter := range parameters {

		kvParameter := strings.Split(parameter, "=")
		updatedParameter := parameter

		if len(sthingsBase.GetAllRegexMatches(kvParameter[1], regexPatternVaultSecretPath)) > 0 {
			secretValue := sthingsCli.GetVaultSecretValue(kvParameter[1], os.Getenv("VAULT_TOKEN"))
			updatedParameter = kvParameter[0] + "=" + secretValue
		}

		updatedParameters = append(updatedParameters, updatedParameter)

	}

	return
}

func VerifyVaultEnvVars() bool {

	if sthingsCli.VerifyEnvVars([]string{"VAULT_ADDR", "VAULT_TOKEN", "VAULT_NAMESPACE"}) {
		fmt.Println("VAULT SET!")
	} else {
		fmt.Println("UNSET!")
	}

	return true
}

func InitalizeTerraform(terraformDir, terraformVersion string) (tf *tfexec.Terraform) {

	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion(terraformVersion)),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		fmt.Println("Error installing Terraform: %s", err)
	}

	tf, err = tfexec.NewTerraform(terraformDir, execPath)
	if err != nil {
		fmt.Println("Error running Terraform: %s", err)
	}

	return

}
