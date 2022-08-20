[![CircleCI](https://circleci.com/gh/SharpSet/sharpdev.svg?style=svg)](https://circleci.com/gh/SharpSet/sharpdev)

![SharpDev](https://files.mcaq.me/039xk.png)
# Development Tools for your projects

Sharpdev aims to make scripts for your projects much easier to use!

# Example Config
```yml
version: 1.0
envfile: .env
values:
  TEST: Values work
  _ROOT: /home/coder/code-server/go/src/sharpdev
  sharpdev_test: /home/coder/code-server/go/src/sharpdev/internal/sharpdev

setup: |
  go build -o _ROOT/internal/sharpdev _ROOT/src
  SETUP="Setup Works"

scripts:
  default: echo "Default works"

  test_echo_single: sharpdev_test echo1
  test_echo_multi: sharpdev_test echo2 Args Works

  test_parent: |
    cd ./internal
    ./sharpdev -p echo3
    cd ..

  test_version: |
    cd ./internal
    echo $(./sharpdev -v) - Version Checks - $(./sharpdev --version)
    cd ..

  test_default: sharpdev_test

  test_env_sub: |
    cd ./tests
    sharpdev_test sub_file
    cd ..

  test_setup: |
    sharpdev_test echo4

  test_skip_setup: |
    if [ "$SETUP" = "Setup Works" ]; then
      echo "Setup Ran: Not Expected\n"
    else
      echo "Setup Failed. Great!\n"
    fi

  test_echo_complicated: |
    sharpdev_test echo5 yay! Duplicates_Work!
    sharpdev_test echo2 Extra Args Works!

  echo1: echo TEST
  echo2: echo "$_ARG1 $_ARG2"
  echo3: echo "Env and Parent ${ECHO:-failed}"
  echo4: echo $SETUP
  echo5: echo "$_ARG1 $_ARG2 $_ARG1"

  full: |
    sharpdev test_echo_single
    sharpdev test_echo_multi
    sharpdev test_echo_complicated
    sharpdev test_parent
    sharpdev test_version
    sharpdev test_default
    sharpdev test_env_sub
    sharpdev test_setup
    sharpdev_test -ss test_skip_setup

```

# Installation
On linux, just run:
```console
sudo curl -s -L https://github.com/SharpSet/sharpdev/releases/download/1.7/install.sh | sudo bash
```

## Command Options

On linux, just run:
```console
$ sharpdev help

This Application lets you run scripts set in your sharpdev.yml file.

Note that if no file is found in the dir you are in, it will instead search in ./env

It Supports:
        - env vars in the form $VAR or ${VAR}
        - Multiline commands with |
        - Inputting Args with env vars like $_ARG{1, 2, 3, 4, etc}

Flags:
        -p Uses a parent sharpdev.yml file

If no script is called, the "default" script will be run.

Here are all the scripts you have available:

test_echo_single || test_version || test_default || test_env_sub || echo1 || echo2 || echo4 || full || test_echo_multi || test_parent || test_setup || echo3 ||
```

## Maintainers

- [Adam McArthur](https://adam.mcaq.me)
