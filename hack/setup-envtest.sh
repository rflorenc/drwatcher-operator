set -eu

# */
# Required binaries to use envtest
# */

K8S_VER=v1.18.2
ETCD_VER=v3.4.3
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/')
ETCD_EXT="tar.gz"

rm -rf ../testbin &> /dev/null
mkdir -p ../testbin

[[ -x ../testbin/etcd ]] || curl -L https://storage.googleapis.com/etcd/${ETCD_VER}/etcd-${ETCD_VER}-${OS}-${ARCH}.${ETCD_EXT} | tar zx -C ../testbin --strip-components=1 etcd-${ETCD_VER}-${OS}-${ARCH}/etcd
[[ -x ../testbin/kube-apiserver && -x ../testbin/kubectl ]] || curl -L https://dl.k8s.io/${K8S_VER}/kubernetes-server-${OS}-${ARCH}.tar.gz | tar zx -C ../testbin --strip-components=3 kubernetes/server/bin/kube-apiserver kubernetes/server/bin/kubectl

