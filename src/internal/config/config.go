package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// use struct tags to define parsing params
type configuration struct {
	Port         int    `parse:"env_name=PORT,default=3000"`           // Port to run the server - defaults to 3000
	DbPath       string `parse:"env_name=DB_PATH,required=true"`       // Path to the DB SQLite file (will be created if not exists)
	WebRootPath  string `parse:"env_name=WEB_ROOT_PATH,required=true"` // Path to the static directory for frontend
	CookieSecret string `parse:"env_name=COOKIE_SECRET,required=true"` // Secret for generating JWT token
	Domain       string `parse:"env_name=DOMAIN"`                      // Optional! Should be set if the auth cookie needs to be present across subdomains (required for reverse proxies with subdomain routing)
	Secure       bool   `parse:"env_name=SECURE_COOKIE,default=true"`  // Optional - defaults true! Set to false to work on localhost (SHOULD NOT USE IN PRODUCTION)
	AuthHost     string `parse:"env_name=AUTH_HOST"`                   // Optional - Set to the login auth domain - redirected here when not authenticated
}

var Config *configuration

func Load() error {
	Config = &configuration{}

	// Get the type of our configuration struct to extract the tags
	t := reflect.TypeOf(*Config)

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get("parse")
		tags := extractTag(tag)

		// Extract from env
		envName := tags["env_name"]
		defaultValue := tags["default"]
		required := tags["required"] == "true" //required should be explicitly set to true

		value := os.Getenv(envName)
		if value == "" {
			if required {
				return fmt.Errorf("required env '%s' not set", envName)
			}
			value = defaultValue
		}

		configField := reflect.ValueOf(Config).Elem().FieldByName(field.Name)

		switch field.Type.Kind() {
		case reflect.String:
			configField.SetString(value)

		// Int values
		case reflect.Int:
			i, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("env '%s' cannot parse value: %s", envName, value)
			}
			configField.SetInt(int64(i))

		// Boolean values (only accept true,TRUE,false,FALSE as valid values)
		case reflect.Bool:
			val := strings.ToLower(value)
			if val == "true" {
				configField.SetBool(true)
			} else if val == "false" {
				configField.SetBool(false)
			} else {
				return fmt.Errorf("env '%s' cannot parse value: %s", envName, value)
			}

		// Struct shouldn't have other types, implement above if needed
		default:
			return fmt.Errorf("unsupported type %v", field.Type.Kind())
		}
	}
	return nil
}

// No validation is done, Ensure the struct tag is valid
func extractTag(tag string) map[string]string {
	// Split the tag value by comma
	parts := strings.Split(tag, ",")

	// Create a map to store the key-value pairs
	m := make(map[string]string)

	// Iterate over all key-value pairs
	for _, part := range parts {
		kv := strings.Split(part, "=")
		m[kv[0]] = kv[1]
	}
	return m
}
