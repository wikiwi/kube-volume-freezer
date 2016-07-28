# kube-volume-freezer
Freeze Pod Volumes in Kubernetes for the purpose of producing a live snapshot.

[![Build Status Widget]][Build Status]
[![GoDoc Widget]][GoDoc]
[![Coverage Status Widget]][Coverage Status]
[![Code Climate Widget]][Code Climate]
[![MicroBadger Version Widget]][MicroBadger Version]

[Build Status]: https://travis-ci.org/wikiwi/kube-volume-freezer
[Build Status Widget]: https://travis-ci.org/wikiwi/kube-volume-freezer.svg?branch=master
[GoDoc]: https://godoc.org/github.com/wikiwi/kube-volume-freezer
[GoDoc Widget]: https://godoc.org/github.com/wikiwi/kube-volume-freezer?status.svg
[Coverage Status]: https://coveralls.io/github/wikiwi/kube-volume-freezer?branch=master
[Coverage Status Widget]: https://coveralls.io/repos/github/wikiwi/kube-volume-freezer/badge.svg?branch=master
[Code Climate]: https://codeclimate.com/github/wikiwi/kube-volume-freezer
[Code Climate Widget]: https://codeclimate.com/github/wikiwi/kube-volume-freezer/badges/gpa.svg
[MicroBadger Version]: http://microbadger.com/#/images/wikiwi/kube-volume-freezer
[MicroBadger Version Widget]: https://images.microbadger.com/badges/version/wikiwi/kube-volume-freezer.svg

## Use-Case
- You want to sync and freeze one or multiple Kubernetes Pod Volumes before creating a live snapshot without adding additional capabilities to your Pods. ([GCE-Guide on creating snapshots](https://cloud.google.com/compute/docs/disks/create-snapshots))

## Architecture
- [kvf-minion](https://github.com/wikiwi/kube-volume-freezer/blob/master/docs/kvf-minion.md) is run on every Node with required privileges and perform the actual syncing, freezing and thawing of local Volumes on the host.
- [kvf-apiserver](https://github.com/wikiwi/kube-volume-freezer/blob/master/docs/kvf-apiserver.md) delegates client requests to the correct Minion.
- [kvfctl](https://github.com/wikiwi/kube-volume-freezer/blob/master/docs/kvfctl.md) is a command-line interface to the [kvf-apiserver](https://github.com/wikiwi/kube-volume-freezer/blob/master/docs/kvf-apiserver.md).

## Install
The folder [manifests](https://github.com/wikiwi/kube-volume-freezer/tree/master/manifests) contains an example deployment of kube-volume-freezer with token protected Minions and API server.

Steps to install:

- Set your base64 encoded tokens in the `kcf-secret.yaml` file.
- Add `kvf-secret.yaml`, `kvf-daemonset.yaml`, `kvf-deployment.yaml` and `kvf-svc.yaml` to your Kubernetes.

## Example
The following example shows how to create a live snapshot from a running system on GCE using `kubectl`, `kvfctl`, and `gcloud`.

    #!/bin/bash
    # Open a local port to the kube-volume-freezer service.
    # This is not needed when running inside a Kubernetes Cluster.
    kubectl port-forward kube-volume-freezer-master-1053963144-7uxa2 8080:8080 &
    PID=$!

    # Freeze Volume named "data" in Pod "gitlab-3323024633-063kf".
    kvfctl freeze --address localhost:8080 --token "my-token" gitlab-3323024633-063kf data

    # Create snapshot on GCE associated with the Pod.
    gcloud compute disks snapshot gitlab-disk --zone europe-west1-b --snapshot-names "gitlab-disk-$(date +"%Y%m%d%H%M%S")"

    # Thaw Volume.
    kvfctl thaw --address localhost:8080 --token "my-token" gitlab-3323024633-063kf data

    # Close local port.
    kill -TERM ${PID}

In a more [complex example](https://gist.github.com/cvle/148440760de3156e7d2394a4c7795193) you can use kubectl to detect persistent disks of different deployments automatically and perform live snapshotting simultaneously.

## Library
  - [golang client](https://godoc.org/github.com/wikiwi/kube-volume-freezer/pkg/client)

