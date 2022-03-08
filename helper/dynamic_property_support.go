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
	"fmt"
	"github.com/buildpacks/libcnb"
	"github.com/magiconair/properties"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/bindings"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"os"
	"path/filepath"
	"strings"
)

const prefix = "BPL_TOMCAT_ENV_"

type DynamicPropertySupport struct {
	Bindings libcnb.Bindings
	Logger   bard.Logger
}

func (a DynamicPropertySupport) Execute() (map[string]string, error) {
	systemProperties := make(map[string]string)

	bindingName := sherpa.GetEnvWithDefault("BPL_TOMCAT_BINDING_NAME", "")
	if bindingName != "" {
		b, ok, err := bindings.ResolveOne(a.Bindings, bindings.OfType(bindingName))
		if err != nil {
			return nil, fmt.Errorf("unable to resolve binding %s\n%w", bindingName, err)
		} else if ok {
			a.Logger.Infof("Configuring properties for binding %s", b.Name)

			for k, v := range b.Secret {
				systemProperties[fmt.Sprintf("%s.%s", b.Name, k)] = v
			}
		}
	}

	// get all environment variables named BPL_TOMCAT_ENV_*
	envVariables := os.Environ()
	for _, envVariable := range envVariables {
		parts := strings.SplitN(envVariable, "=", 2)
		if strings.HasPrefix(parts[0], prefix) {
			a.Logger.Debugf("Configuring property from env var %s", parts[0])
			systemProperties[toSystemProperty(parts[0])] = parts[1]
		}
	}

	catalinaBase, err := sherpa.GetEnvRequired("CATALINA_BASE")
	if err != nil {
		return nil, err
	}

	catalinaProperties := filepath.Join(catalinaBase, "conf", "catalina.properties")

	err = a.writeToCatalinaProperties(catalinaProperties, systemProperties)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (a DynamicPropertySupport) writeToCatalinaProperties(catalinaProperties string, systemProperties map[string]string) error {
	p := properties.NewProperties()

	a.Logger.Infof("Writing properties to %s", catalinaProperties)
	f, err := os.OpenFile(catalinaProperties, os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	for k, v := range systemProperties {
		p.Set(k, v)
	}

	_, err = p.Write(f, properties.UTF8)
	if err != nil {
		return err
	}

	return nil
}

func toSystemProperty(env string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(env, prefix)), "_", ".")
}
