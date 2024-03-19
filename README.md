# binary-version-switcher

## Install

```shell
brew install maahsome/tap/binary-version-switcher --formula
```

## Use / TLDR

_NOTE:_ currently the `v0.0.1` artifacts only support MacOS

TODO: build a pipeline to build everything

```shell
# initial run, just run the application
# select the locations
binary-version-switcher
# download the plugins
binary-version-switcher plugins get
# list a version
binary-version-switcher terraform versions
# activate a terraform version
binary-version-swticher terraform activate -v "1.6.6"
# test the version
terraform version
```

## Description

A quick way to switch between binary versions of common tools.

Some tools that I use `helm`, `kubectl`, and `terraform` generally require a
specific version when interacting with certain `targets/sources` and using `brew`
isn't always possible as most `brew formula` are one-way upgrades, lacking a way
to downgrade to an older version.

So, this tool will download different versions of an application into a `binPath`
from the `config.yaml` file, storing the versions in
`<binPath>/<appName>/<semver>/<appName>`, then will create a symbolic link in the
path defined in `symLinkPath` in the `config.yaml` file.  As long as the `path`
defined in `symLinkPath` is before your normal path where `brew` installs are
done, or you define `symLinkPath` to be the same as your `brew` bin path, the
specific version will be used.

I tend to use `/usr/local/bvs` as my `binPath`, and `/usr/local/bin` as my
`symLinkPath`, though on a multi-user system, say a hosted VM in the cloud that
many team members access, that might not be great.  For a mult-user system
there are some defaults that will be presented:

`binPath` = `<HOMEDIR>/.config/binary-version-switcher/bvs`
`symLinkPath` = `<HOMEDIR>/.config/binary-version-switcher/bin`

You'll just need to add the `symLinkPath` to the front of your `PATH` environment
variable in your shell startup script.

The different applications are build as `golang` `plugin` binaries.  Running
`binary-version-switcher plugin get` will prompt for a list of supported
`applications`.  The sources for the plugins can be found [here](https://github.com/maahsome/binary-version-switcher-plugins).

## Design

This project isn't just because I could use a quick away to switch between versions
of binary applications.  It was for learning some things that I hadn't used in
golang yet.  In this case, the `plugin` functionality.

My experience was the same as others, the oddity of trying to keep the `go.mod`
references in sync to ensure they are the same as `this` application was quite
annoying.  And clearly creating an interface of `apps` defining the functions
would have been _MUCH_ cleaner and easier to add support for different apps and
all that jazz.  In the end, it was a nice journey and it works, so I'll leave it
this way.

## Refs

- [kubectl plugins](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/)
- [Cobra and plugin](https://blog.chmouel.com/2022/03/18/show-plugins-list-to-a-cli-when-using-gos-cobra-library/)
  - [code](https://github.com/tektoncd/cli/pull/1535/files)
- [Hashicorp gRPC Plugin](https://github.com/hashicorp/go-plugin)

