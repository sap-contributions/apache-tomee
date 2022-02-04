#!/bin/bash

set -euo pipefail

git stash

grep -Rl "paketo-buildpacks/apache-tomcat" * .github | grep -v migrate | xargs sed -i 's|paketo-buildpacks/apache-tomcat|garethjevans/apache-tomcat|'
grep -Rl "gcr.io/garethjevans" * .github | grep -v migrate | xargs sed -i 's|gcr.io/garethjevans|ghcr.io/garethjevans|'

rm .github/CODEOWNERS

echo "golang 1.17.5" > .tool-versions

sed -i 's|jobs:|\nenv:\n\njobs:|' .github/workflows/create-package.yml
sed -i 's|ubuntu-latest|ubuntu-latest\n        permissions:\n|' .github/workflows/create-package.yml

echo "Adding env"
yq e -i '.env.REGISTRY = "ghcr.io"' .github/workflows/create-package.yml
yq e -i '.env.IMAGE_NAME = "${{ github.repository }}"' .github/workflows/create-package.yml

echo "Adding permissions"
yq e -i '.jobs.create-package.permissions.contents = "read"' .github/workflows/create-package.yml
yq e -i '.jobs.create-package.permissions.packages = "write"' .github/workflows/create-package.yml
 
sed -i 's/Docker login gcr.io/Docker login ghcr.io/' .github/workflows/create-package.yml

sed -i 's/password: \${{ secrets.JAVA_GCLOUD_SERVICE_ACCOUNT_KEY }}/password: \${{ secrets.GITHUB_TOKEN }}/' .github/workflows/create-package.yml
sed -i 's/registry: gcr.io/registry: \${{ env.REGISTRY }}/' .github/workflows/create-package.yml
sed -i 's/username: _json_key/username: \${{ github.actor }}/' .github/workflows/create-package.yml
sed -i 's/if: \${{ true }}/if: \${{ false }}/' .github/workflows/create-package.yml

sed -i 's/0-9/0-9]+\\.[0-9/' .github/workflows/create-package.yml

go mod tidy

go build ./...
go test ./...

echo "Add the secret JAVA_GITHUB_TOKEN"
gh secret set JAVA_GITHUB_TOKEN
