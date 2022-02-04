package helper

import (
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"os"
	"strings"
)

// if enabled, add properties to define the datasource
type ManagedDataSourceSupport struct {
	Logger bard.Logger
}

func (a ManagedDataSourceSupport) Execute() (map[string]string, error) {
	a.Logger.Info("################################################")
	a.Logger.Info("checking ManagedDataSourceSupport...")
	a.Logger.Info("################################################")

	if _, ok := os.LookupEnv("BPL_TOMCAT_MANAGED_DATASOURCE_ENABLED"); !ok {
		return nil, nil
	}

	a.Logger.Info("Tomcat Managed Datasource Enabled")

	var values []string
	if s, ok := os.LookupEnv("JAVA_TOOL_OPTIONS"); ok {
		values = append(values, s)
	}

	values = append(values, "-Dmanaged.datasource.name="+sherpa.GetEnvWithDefault("BPL_TOMCAT_MANAGED_DATASOURCE_NAME", "jdbc/datasource"))
	values = append(values, "-Dmanaged.datasource.auth="+sherpa.GetEnvWithDefault("BPL_TOMCAT_MANAGED_DATASOURCE_AUTH", "Container"))
	values = append(values, "-Dmanaged.datasource.type="+sherpa.GetEnvWithDefault("BPL_TOMCAT_MANAGED_DATASOURCE_TYPE", "javax.sql.DataSource"))

	driver, err := sherpa.GetEnvRequired("BPL_TOMCAT_MANAGED_DATASOURCE_DRIVER")
	if err != nil {
		return nil, err
	}

	values = append(values, "-Dmanaged.datasource.driver="+driver)

	username, err := sherpa.GetEnvRequired("BPL_TOMCAT_MANAGED_DATASOURCE_USERNAME")
	if err != nil {
		return nil, err
	}

	values = append(values, "-Dmanaged.datasource.username="+username)

	password, err := sherpa.GetEnvRequired("BPL_TOMCAT_MANAGED_DATASOURCE_PASSWORD")
	if err != nil {
		return nil, err
	}

	values = append(values, "-Dmanaged.datasource.password="+password)

	url, err := sherpa.GetEnvRequired("BPL_TOMCAT_MANAGED_DATASOURCE_URL")
	if err != nil {
		return nil, err
	}

	values = append(values, "-Dmanaged.datasource.url="+url)

	values = append(values, "-Dmanaged.datasource.maxtotal="+sherpa.GetEnvWithDefault("BPL_TOMCAT_MANAGED_DATASOURCE_MAX_TOTAL", "20"))
	values = append(values, "-Dmanaged.datasource.maxidle="+sherpa.GetEnvWithDefault("BPL_TOMCAT_MANAGED_DATASOURCE_MAX_IDLE", "10"))
	values = append(values, "-Dmanaged.datasource.maxwaitmillis="+sherpa.GetEnvWithDefault("BPL_TOMCAT_MANAGED_DATASOURCE_MAX_WAIT_MILLIS", "-1"))

	return map[string]string{"JAVA_TOOL_OPTIONS": strings.Join(values, " ")}, nil
}
