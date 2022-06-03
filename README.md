[![CircleCI](https://circleci.com/gh/Sharpz7/sharpdev.svg?style=svg)](https://circleci.com/gh/Sharpz7/sharpdev)


# SharpDev | Development Tools for your projects

Sharpdev aims to make scripts for your projects much easier to use!

# Example Config
```yml
version: 1.0
envfile: .env
values:
  TEST: Hello World!

scripts:
  build: go build -o ./internal/sharpdev ./src
  test1: sharpdev build && ./internal/sharpdev echo1
  test2: sharpdev build && ./internal/sharpdev echo2 $_ARG1
  test3: |
    sharpdev build
    cd ./internal
    ./sharpdev -p echo1

  echo1: echo TEST
  echo2: echo $_ARG1

  full: |
    sharpdev test1
    sharpdev test2 Test2
    sharpdev test3
```

# Installation
On linux, just run:
```console
sudo curl -s -L https://github.com/Sharpz7/sharpdev/releases/download/1.3/install.sh | sudo bash
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

echo2 || full || build || test1 || test2 || test3 || echo1 ||
```

## Maintainers

- [Adam McArthur](https://adam.mcaq.me)
