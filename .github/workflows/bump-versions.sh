#!/bin/bash

set +x

if [ -z "${PREFIX}" ]; then
	echo "Expecting \$PREFIX to be set"
	exit 1
fi
if [ -z "${GITHUB_REF_NAME}" ]; then
	echo "Expecting \$GITHUB_REF_NAME to be set"
	exit 1
fi
if [ -z "${GITHUB_RUN_NUMBER}" ]; then
	echo "Expecting \$GITHUB_RUN_NUMBER to be set"
	exit 1
fi
if [ -z "${ENVIRONMENT}" ]; then
	echo "Expecting \$ENVIRONMENT to be set"
	exit 1
fi

cd ./src
go run ./cmd/bump-versions -module-namespace "${PREFIX}" -provider-namespace "${PREFIX}"
cd ../

BRANCH="${GITHUB_REF_NAME}-bump-versions-${GITHUB_RUN_NUMBER}"
if [ "${ENVIRONMENT}" = "Production" ]; then
	export BRANCH="${GITHUB_REF_NAME}"
fi

# Try to get to the latest HEAD
git fetch origin $BRANCH
git checkout -b $BRANCH

# Commit changes
git status
git add ./modules ./providers
git config user.email "core@opentofu.org"
git config user.name "OpenTofu Core Development Team"
git commit -m "Automated bump of versions for providers and modules (${PREFIX})"

# Racing with other jobs, try a few times to push changes
for i in {0..30}; do
	git push -u origin $BRANCH
	if [ $? = 0 ]; then
		echo "Providers and modules changes were pushed to branch: ${BRANCH}"
		exit 0
	fi
	sleep 1
	git fetch origin
	git rebase origin/$BRANCH
done

echo "Failed to push providers and modules changes to branch: ${BRANCH}"
exit -1

