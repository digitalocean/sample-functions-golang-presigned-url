# Sample Function: Go "Presigned URL"

## Introduction

This repository contains a sample presigned URL function written in Go. You are able to choose to get a presigned URL to upload a file to a DigitalOcean Space or to download a file from a DigitalOcean Space. You can deploy it on DigitalOcean's App Platform as a Serverless Function component.
Documentation is available at https://docs.digitalocean.com/products/functions.

### Requirements

* You need a DigitalOcean account. If you don't already have one, you can sign up at [https://cloud.digitalocean.com/registrations/new](https://cloud.digitalocean.com/registrations/new).
* You need a DigitalOcean Space. If you don't have one, you can create one at https://www.digitalocean.com/products/spaces.
* You need to add your `SPACES_KEY`, `SPACES_SECRET`, `BUCKET`, and `REGION` to the `.env` file to connect to Spaces API as well as your bucket.
* To deploy from the command line, you will need the [DigitalOcean `doctl` CLI](https://github.com/digitalocean/doctl/releases).

## Deploying the Function

```bash
# clone this repo
git clone https://github.com/digitalocean/sample-functions-golang-presigned-url
```

```
# deploy the project, using a remote build so that compiled executable matched runtime environment
> doctl serverless deploy sample-functions-golang-presigned-url --remote-build
Deploying 'sample-functions-golang-presigned-url'
  to namespace 'fn-...'
  on host 'https://faas-...'
Submitted action 'url' for remote building and deployment in runtime go:default (id: ...)

Deployed functions ('doctl sls fn get <funcName> --url' for URL):
  - presign/url
```

## Using the Function

```bash
doctl serverless functions invoke presign/url -p filename:new-file.txt type:GET
```
```json
{
  "body": "{presigned url}"
}
```

### To get a presigned url using curl:
```
curl -X PUT -H 'Content-Type: application/json' {your-DO-app-url} -d '{"filename":"{filename}", "type":"GET or PUT"}'
```

### To Upload or Download a file using curl:
```
curl -X PUT -d 'The contents of the file.' "{url}"
```

### Learn More

You can learn more about Functions and App Platform integration in [the official App Platform Documentation](https://www.digitalocean.com/docs/app-platform/).
