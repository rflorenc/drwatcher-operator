[![Go Report Card](https://goreportcard.com/badge/github.com/rflorenc/drwatcher-operator)](https://goreportcard.com/report/github.com/rflorenc/drwatcher-operator)

# DRWatcher Self Service Operator

## Summary

The DR (Disaster Restic|Recovery) Watcher is used for enabling Velero based Self Service backups and schedules in OpenShift v4.x, Kubernetes v1.13 or later.
An example backup infra can be setup by following the instructions here: [https://github.com/rflorenc/openshift-backup-infra](https://github.com/rflorenc/openshift-backup-infra)

## DR Watcher Operator Functionality

1. Pre-requisites to starting DRWatcher:
    + A fully functioning Velero installation, including a default BackupStorageLocation, for example Noobaa.
    + The DRWatcher CustomResourceDefinition (config/crd/bases) must be applied.
    + A DRWatcher CustomResource (config/samples/) must be installed in the namespace in which we want to self service backup.

2. A DRWatcher Custom Resource must exist for each namespace to self service:
    + When `readyForBackup: true` and a `schedule` is defined, DRWatcher will create a Scheduled backup for the namespace.
    + When `readyForBackup: true` and a `schedule` is absent, DRWatcher will create an immediate Backup for the namespace.
    + When `readyForBackup: false`, the reconciler only logs the existing Restic annotations (`backup.velero.io/backup-volumes`)

See `config/samples`:

```yaml
apiVersion: dr.seven/v1
kind: DRWatcher
metadata:
  name: drwatcher-sample
spec:
  readyForBackup: true
  schedule: '0 1 * * *'
```

## Requirements

+ Access to OpenShift v4.1, Kubernetes v1.13 or later.
+ A working Velero/Konveyor/Restic/Noobaa installation.

## Install via Helm (in-cluster)

```shell
helm-wrapper.sh -u
```

The operator will watch for changes done in `DRWatcher` instances created in the target namespaces.
See `config/samples` for example CRs to use.

## How to build and run locally (off-cluster)

```shell
make
build/bin/drwatcher-operator
```

## How to generate manifests

```shell
make manifests
```

## How to run unit tests

Launch the tests locally by running

```shell
make test
```
