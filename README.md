# stuttgart-things/sthings-snake

play classic mobile game snake in your terminal & create some chaos on a k8s cluster by homerun messages.
Based on [chaossnake](https://github.com/deggja/chaossnake) - @deggja thank you for your great work.

## USAGE

```bash
export HOMERUN_TOKEN=""
export HOMERUN_ADDR="https://homerun.homerun-dev.."

sthings-snake
```

## DEPLOYMENT

<details><summary>BY GITHUB RELEASE</summary>

```bash
VERSION=v1.1.0
BIN_DIR=/usr/bin
cd /tmp && wget https://github.com/stuttgart-things/sthings-snake/releases/download/${VERSION}/sthings-snake_Linux_x86_64.tar.gz
tar xvfz sthings-snake_Linux_x86_64.tar.gz
sudo mv sthings-snake ${BIN_DIR}/sthings-snake
sudo chmod +x ${BIN_DIR}/sthings-snake
rm -rf CHANGELOG.md README.md LICENSE sthings-snake_Linux_x86_64.tar.gz
cd -
```

</details>

## DEV

<details><summary>CREATE ENV FILE</summary>

.env file needed for Taskfile

```bash
cat <<EOF > .env
HOMERUN_TOKEN=""
HOMERUN_ADDR="https://homerun.homerun-dev..."
EOF
```

</details>


```bash
task: Available tasks for this project:
* branch:           Create branch from main
* build:            Build the binary
* check:            Run pre-commit hooks
* commit:           Commit + push code into branch
* goreleaser:       Release bins w/ goreleaser
* install:          Installs Tetrigo        (aliases: i)
* lint:             Runs golangci-lint      (aliases: l)
* pr:               Create pull request into main
* release:          Push new version
* run:              Run the Go project
* test:             Runs test suite      (aliases: t)
* lint:fix:         Runs golangci-lint and fixes any issues
```
