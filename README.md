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
  revert: git revert $_ARG1..HEAD
  list: git branch
  archive: |
   git tag archive/$_ARG1 $_ARG1 &&
   git branch -d $_ARG1
  test1: sharpdev build && ./internal/sharpdev $_ARG1
  test2: sharpdev build && ./internal/sharpdev $_ARG1 $_ARG2
  new1: |
    echo "Hello World!"
    echo "Test 2"

  new2: echo $_ARG1

  full: |
    sharpdev test1 help &&
    sharpdev test2 new2 TEXT &&
    sharpdev test1 new1
```

# Installation
On linux, just run:
```console
sudo curl -s -L https://github.com/Sharpz7/sharpdev/releases/download/1.0/install.sh | sudo bash
```

## Command Options

On linux, just run:
```console
sharpdev --help

This Application lets you run scripts set in your sharpdev.yml file.

It Supports:
        - env vars in the form $VAR or ${VAR}
        - Multiline commands with |
        - Inputting Args with env vars like $@ARG{1, 2, 3, 4, etc}

Here are all the scripts you have available:

revert || archive || new2 || new1 || full || build || list || test1 || test2 ||
```

## Maintainers

- [Adam McArthur](https://adam.mcaq.me)
