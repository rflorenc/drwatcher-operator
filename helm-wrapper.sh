#!/usr/bin/env bash

dryrun=""
helm_chart_name=${PWD##*/}

function print_help() {
cat << EOF
    Usage: ./helm_wrapper.sh
    -u          helm upgrade <namespace>
    -d          helm dry-run debug
    -r          remove release
EOF
  exit 1
}

function upgrade_helm() {
  namespace=${1:-prometheusrule-webhook}
  helm upgrade ${helm_chart_name} \
    --namespace ${namespace} \
    --install \
    --debug \
    --create-namespace \
    $dryrun \
    ./helm/${helm_chart_name} \
    -f ./helm/${helm_chart_name}/values.yaml

  echo -e "\nInstalled in namespace: ${namespace}"
}

function remove_helm() {
  helm uninstall ${helm_chart_name} --debug
}


case $1 in
  "-h") print_help ;;
  "-u") upgrade_helm $2;;
  "-d") dryrun="--dry-run"; upgrade_helm $2;;
  "-r") remove_helm ;;
  *) print_help ;;
esac
