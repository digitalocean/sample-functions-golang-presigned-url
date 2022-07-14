# Sample Function: Go "Presigned URL"

## Introduction

This repository contains a sample presigned URL function written in Go. You are able to choose to get a presigned URL to upload a file to a DigitalOcean Space or to download a file from a DigitalOcean Space.

### To get a url:
curl -X PUT -H 'Content-Type: application/json' {your-DO-app-url} -d '{"filename":"{filename}", "type":"GET or PUT"}'

### To Upload or Download the file:
curl -X PUT -d 'The contents of the file.' "{url}"