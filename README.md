# TASC

A generic *T*ool for *A*ssembling *S*ource *C*ode.

## About

Sometime you just need a generic tool to fetch some source code from various
locations and assemble it in a way you specify. My goal for this project is to
create a simple reproducable method for assembling code.

Tasc assembles your source code based on a YAML formatted "manifest" which lists
which projects should go where.

## Dependendies

[Git](https://git-scm.com/downloads "Git downloads") (if you want to fetch code
with git) and [GNU Patch](http://www.gnu.org/s/patch/ "GNU Patch project page")
(if you want to apply patches).

## Installation

Grab a binary from the releases. Or if you have Go:

```
$ go get github.com/HCPSS/tasc
```

## Usage

```
$ tasc -h
  _______        _____  _____
 |__   __|/\    / ____|/ ____|
    | |  /  \  | (___ | |
    | | / /\ \  \___ \| |
    | |/ ____ \ ____) | |____
    |_/_/    \_\_____/ \_____|

A tool of assembling source code.

Version v0.1.0
Copyright (C) 2016 Howard County Public Schools
Distributed under the terms of the MIT license
Written by Brendan Anderson

  -destination string
    	Where to build the project (default "./")
  -manifest string
    	Name of the manifest file. (default "manifest.yml")
  -params string
    	A JSON encoded string with extra parameters. (default "{}")
  -v	Print the version.
  -version
    	Print the version.
```

The *destination* is pretty self explanitory. The *manifest* and *params* are
discussed below.

## Manifest

You need to tell tasc how to assemble your project witha YAML formatted
"manifest". The manifest lists projects and patches. For example:

```yaml
---
# Code to fetch and assemble.
projects:
  -
    # How are we going to use to fetch the code?
    # Options are git, zip, and local.
    provider: git

    # For git, we want the http URI of the git repo.
    source: "https://github.com/moodle/moodle.git"

    # Version can be a branch, a tag (in the format tags/<tag-name>), or a
    # commit hash, as shown here.
    version: badfcb70e4e59ca0a3d4fc29b34174eb06f89b95

    # Tags a extra metadata that tell tasc how to process the project.
    tags:
      # Projects are downloaded simultaniously unless they have the tag
      # "blocking". All blocking projects are processed before other projects.
      # In this case, Moodle is our root project so it will set up the folder
      # structure that other projects will use.
      - blocking

      # Sticky makes sure that the project is shown at the top of the processing
      # list.
      - sticky


    # In this example, we construct a basic oauth request URI so that we can
    # access private repos. We will provide a value for the github_access_token
    # placeholder with the params flag when we run tasc.
    - provider: git

      # Where to download the project to, relative to the root destination. If
      # no destination is provided, the project is placed in the project root.
      destination: local

      # Should we rename the directory? This project will be located at:
      # <project-root>/local/provisioner
      rename: provisioner
      source: "https://{github_access_token}:x-oauth-basic@github.com/HCPSS/moodle-enrol_mandatory.git"
      version: tags/v2.0.0

    # Example zip provider:
    - provider: zip
      source: "https://moodle.org/plugins/download.php/8086/format_grid_moodle28_2015022500.zip"
      destination: course/format

    # The local provider simply gets files from the filesystem.
    - provider: local
      source: "{manifest_dir}/customfiles"
      destination: "{destination_dir}/custom"

# Should we perform any patches once the code is assembled?
patches:
  # Patch the forum to add debugging during cron and modify the template
  -
    # Currently the only patch type that is supported is patch_file. This method
    # patches a single file with the GNU patch program.
    type:         patch_file

    # {manifest_dir} and {destination_dir} are replaced with the manifest
    # directory and the destination direcory at runtime. So you can use those
    # without specifying them in params.
    source:       "{manifest_dir}/patches/mod_forum_lib.php.patch"
    destination:  "{destination_dir}/mod/forum/lib.php"
```

## Params

Params are a JSON encoded list of extra parameters to pass to tasc. These
parameters will be found (in the format {param_name}) and replaced in the
manifest once the file is parsed.

Using the above manifest as an example, there is a placeholder in the file for
github_access_token. I pass in the value for the placeholder like this:

```
$ tasc -params='{"github_access_token": "MYTOKEN"}'
```

Notice that there are also placeholders for manifest_dir and destination_dir.
These are added to the parameters automatically so there is no need to specify
them.

## License

Distributed under the terms of the MIT license.
