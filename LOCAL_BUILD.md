# Local Build

Use this set of commands to perform a local build for tesing.

```bash
SEMVER=v0.0.1; echo ${SEMVER}
BUILD_DATE=$(gdate --utc +%FT%T.%3NZ); echo ${BUILD_DATE}
GIT_COMMIT=$(git rev-parse HEAD); echo ${GIT_COMMIT}

MODULE_NAME=binary-version-switcher
go build -ldflags "-X ${MODULE_NAME}/cmd.semVer=${SEMVER} -X ${MODULE_NAME}/cmd.buildDate=${BUILD_DATE} -X ${MODULE_NAME}/cmd.gitCommit=${GIT_COMMIT} -X ${MODULE_NAME}/cmd.gitRef=/refs/tags/${SEMVER}" && \
./binary-version-switcher version | jq .

if [[ -d ~/tbin ]]; then
  cp ./binary-version-switcher ~/tbin/binary-version-switcher
fi
```

