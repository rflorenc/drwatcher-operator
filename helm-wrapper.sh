#!/usr/bin/env bash

dryrun=""
helm_chart_name=drwatcher-operator
namespace=example

function print_help() {
cat << EOF
    Usage: ./helm_wrapper.sh
    -u          helm upgrade
    -d          helm dry-run debug
    -r          remove release
EOF
  exit 1
}

function upgrade_helm() {
  helm upgrade ${helm_chart_name} \
    --namespace ${namespace} \
    --install \
    --debug \
    $dryrun \
    ./helm/${helm_chart_name} \
    -f ./helm/${helm_chart_name}/values.yaml
}

function remove_helm() {
  helm uninstall drwatcher-operator --debug
}


case $1 in
  "-h") print_help ;;
  "-u") upgrade_helm ;;
  "-d") dryrun="--dry-run"; upgrade_helm ;;
  "-r") remove_helm ;;
  *) print_help ;;
esac
