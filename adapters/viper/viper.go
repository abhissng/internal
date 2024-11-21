package viper

import (
	"fmt"
	"os"
	"regexp"

	"github.com/abhissng/internal/vault"
	"github.com/spf13/viper"
)

type Viper struct {
	configName string
	configType string
	configPath string // it should only contain the absolute path for the folder rest other details will be added by sdk
}

func NewViper(configName, configType, configPath string) *Viper {
	env := os.Getenv("RunMode")
	return &Viper{
		configName: configName,
		configType: configType,
		configPath: configPath + "/" + env + "/",
	}

}

func (v *Viper) InitialiseViper() error {
	viper.SetConfigName(v.configName) // Name of configuration file
	viper.SetConfigType(v.configType) // Configuration file type
	viper.AddConfigPath(v.configPath) // Look for configuration file in the given directory

	// Enable Viper to read environment variables
	viper.AutomaticEnv()

	// Attempt to read configuration file
	if err := viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("error reading configuration file: %s", err)
		return err
	}
	if err := loadAndReplaceConfig(); err != nil {
		fmt.Printf("Error loading and replacing config: %v", err)
		return err
	}
	return nil
}

// Function to load configuration and replace placeholders with values fetched from Vault
func loadAndReplaceConfig() error {
	vlt := vault.NewVault(os.Getenv("RunMode"), viper.GetString("PROJECT_ID"), viper.GetString("VAULT_PATH"))
	vlt.InitVaultClient()

	//  Iterate through all settings and replace placeholders
	for key, value := range viper.AllSettings() {
		// Check if the value is a string and contains placeholders like {{.ENV.DBPASSWORD}}
		if strValue, ok := value.(string); ok {
			// Replace placeholders in the string
			updatedValue, err := replacePlaceholdersWithVault(strValue, vlt)
			if err != nil {
				fmt.Printf("Error fetching secret from Vault for key %s: %v", key, err)
				continue
			}
			// Update the Viper configuration with the new value
			viper.Set(key, updatedValue)
		}
	}
	return nil
}

// Function to replace placeholders with values from Vault
func replacePlaceholdersWithVault(configContent string, vault *vault.Vault) (string, error) {
	// Define the regular expression to match placeholders like {{.env.DBPASSWORD}}
	re := regexp.MustCompile(`\{\{\.(ENV\.[A-Za-z0-9_]+)\}\}`)

	// Replace the placeholders in the config string with the corresponding secret from Vault
	updatedConfig := re.ReplaceAllStringFunc(configContent, func(placeholder string) string {
		// Trim the {{. and }} from the placeholder to get the key directly
		key := placeholder[3 : len(placeholder)-2]

		// Fetch the value from Vault (you can fetch other secrets depending on your Vault setup)
		value, err := vault.FetchVaultValue(key)
		if err != nil {
			fmt.Printf("Error fetching secret %s from Vault: %v", key, err)
			return "" // Return empty string if fetching fails
		}

		return value
	})

	return updatedConfig, nil
}
