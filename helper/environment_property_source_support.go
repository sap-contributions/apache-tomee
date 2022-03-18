/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helper

import (
	"os"
	"strings"
	"github.com/paketo-buildpacks/libpak/bard"
)

type EnvironmentPropertySourceSupport struct {
	Logger bard.Logger
}

func (e EnvironmentPropertySourceSupport) Execute() (map[string]string, error) {
	if _, ok := os.LookupEnv("BPL_TOMEE_ENVIRONMENT_PROPERTY_SUPPORT_ENABLED"); !ok {
		return nil, nil
	}

	e.Logger.Info("Tomee Environment Property Support Enabled")

	var values []string
	if s, ok := os.LookupEnv("JAVA_TOOL_OPTIONS"); ok {
		values = append(values, s)
	}

	values = append(values, "-Dorg.apache.tomcat.util.digester.PROPERTY_SOURCE=org.apache.tomcat.util.digester.EnvironmentPropertySource")

	return map[string]string{"JAVA_TOOL_OPTIONS": strings.Join(values, " ")}, nil
}