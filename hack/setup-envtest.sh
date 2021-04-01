#!/bin/sh

# */
# Required binaries to use controller-runtime envtest
# */


K8S_VER=v1.18.2
ETCD_VER=v3.2.32
OS=linux
ARCH=amd64
EXT="tar.gz"
OUTPUT_BINDIR=/tmp/testbin
GITHUB_URL=https://github.com/etcd-io/etcd/releases/download
K8S_URL=https://dl.k8s.io


rm -rf ${OUTPUT_BINDIR} &> /dev/null
mkdir -p ${OUTPUT_BINDIR}

echo get etcd
curl -s -LO ${GITHUB_URL}/${ETCD_VER}/etcd-${ETCD_VER}-${OS}-${ARCH}.${EXT}
tar -xzf etcd-${ETCD_VER}-${OS}-${ARCH}.${EXT} -C /tmp
cp -v /tmp/etcd-${ETCD_VER}-${OS}-${ARCH}/etcd ${OUTPUT_BINDIR}/

echo get kubectl and kube-apiserver
curl -s -LO ${K8S_URL}/${K8S_VER}/kubernetes-server-${OS}-${ARCH}.${EXT}
tar -xzf kubernetes-server-${OS}-${ARCH}.${EXT} -C /tmp
cp -v /tmp/kubernetes/server/bin/kube-apiserver ${OUTPUT_BINDIR}/
cp -v /tmp/kubernetes/server/bin/kubectl ${OUTPUT_BINDIR}/
