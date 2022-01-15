[![Build](https://github.com/nao1215/gal/actions/workflows/build.yml/badge.svg)](https://github.com/nao1215/gal/actions/workflows/build.yml)

[[日本語](./doc/README.ja.md)]
# gal - generate authors file from git log
gal command generate AUTHORS.md file at current directory. gal command gets the author name and email address from the information in the git log. The information is written in alphabetical order in AUTHORS.md.
# How to install
## Step.1 Install golang
If you don't install golang in your system, please install Golang first. Check the [Go official website](https://go.dev/doc/install) for how to install golang.

## Step2. Install gal
```
$ go install github.com/nao1215/gal/cmd/gal@latest
```

# How to use
You execute the gal command in the directory where .git directory exists. **Please note that the existing AUTHORS.md will be overwritten.**
```
$ gal

$ cat AUTHORS.md 
# Authors List (in alphabetical order)
CHIKAMATSU Naohiro<n.chika156@gmail.com>
TEST User<test@gmail.com>
```

# Contact
If you would like to send comments such as "find a bug" or "request for additional features" to the developer, please use one of the following contacts.

- [GitHub Issue](https://github.com/nao1215/gal/issues)

# LICENSE
The gal project is licensed under the terms of [the Apache License 2.0](./LICENSE).
