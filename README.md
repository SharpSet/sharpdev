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

setup: |
  echo "Setup"

scripts:
  build: go build -o ./internal/sharpdev ./src
  test1: sharpdev build && ./internal/sharpdev echo1
  test2: sharpdev build && ./internal/sharpdev echo2 $_ARG1 $_ARG2

  test3: |
    sharpdev build
    cd ./internal
    ./sharpdev -p echo3
    cd ..

  test4: |
    sharpdev build
    cd ./internal
    ./sharpdev -v
    ./sharpdev --version
    cd ..

  test5: sharpdev build && ./internal/sharpdev

  echo1: echo TEST
  echo2: echo $_ARG1 $_ARG2
  echo3: echo $ECHO

  full: |
    sharpdev test1
    sharpdev test2 Args Works
    sharpdev test3
    sharpdev test4
    sharpdev test5
```

# Installation
On linux, just run:
```console
sudo curl -s -L https://github.com/SharpSet/sharpdev/releases/download/1.5/install.sh | sudo bash
```

## Command Options

On linux, just run:
```console
sharpdev --help

This Application lets you run scripts set in your sharpdev.yml file.

It Supports:
        - env vars in the form $VAR or ${VAR}
        - Multiline commands with |
        - Inputting Args with env vars like $_ARG{1, 2, 3, 4, etc}

Here are all the scripts you have available:

If no script is called, the "default" script will be run.

echo2 || full || build || test1 || test2 || test3 || echo1 ||
```

## Maintainers

- [Adam McArthur](https://adam.mcaq.me)
