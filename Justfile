build: build-vite build-backend

set unstable := true
set dotenv-load := true

REPO := env('REPO', 'UNDEFINED')
NAME := env('NAME', 'UNDEFINED')

BRANCH_DEFAULT := shell("git rev-parse --abbrev-ref HEAD")
GITHUB_REF_NAME := env('GITHUB_REF_NAME', BRANCH_DEFAULT)
GITHUB_HEAD_REF := env('GITHUB_HEAD_REF', '') || GITHUB_REF_NAME
VERSION := trim(read("VERSION.txt"))
BRANCH := replace(GITHUB_HEAD_REF, "/", "-")
SHORT_SHA := trim(shell("git rev-parse --short HEAD"))
TIMESTAMP := datetime("%Y%m%d%H%M%S")
TAG := VERSION + "-" + BRANCH + "-" + SHORT_SHA + "-" + TIMESTAMP

build-vite:
    bun run build

build-backend:
    go build -o ./tmp/main .

build-backend-amd64:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o ./tmp/main .

build-amd64: build-vite build-backend-amd64

docker-build:
    docker build -t go-svelte .

docker-build-amd64:
    docker build --platform linux/amd64 -t go-svelte .

docker-run:
    docker run --rm -it -p 8000:8000 go-svelte

docker-tag-push:
    docker tag go-svelte:latest {{ REPO }}/{{ NAME }}:{{ TAG }}
    docker push {{ REPO }}/{{ NAME }}:{{ TAG }}

docker-release: docker-create-tag build-amd64 docker-build-amd64

docker-publish: docker-tag-push

docker-create-tag:
    @echo "TAG: {{ TAG }}"
    @echo {{ TAG }} > TAG.txt

cr-update:
    gcloud run deploy {{ NAME }} \
    --image ${REPO}/{{ NAME }}:latest \
    --region ${REGION} --project=${PROJECT}

prerequisites:
    bun i

clean:
    -rm -rf tmp node_modules dist
