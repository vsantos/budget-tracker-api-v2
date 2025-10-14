# Running the code within a local Kubernetes cluster

In this section, let's explore how you can run the platform locally using a Kubernetes cluster.

There are many ways of creating a local kubernetes cluster: `kind`, `minikube`, `k3s`, etc. For the sake of simplity, we recommend the usage of `kind`, the only requisite is to have a container engine installed; since `kind` creates a kubernetes cluster by bootstraping containers.

## Create a local k8s cluster

By running the following `kind` command, you will create a functional k8s cluster named `budget-tracker-cluster`:

=== "Shell"
  ```shell
  kind create cluster --name budget-tracker-cluster
  ```
=== "Shell output"
  ```shell
  Creating cluster "budget-tracker-cluster" ...
  ✓ Ensuring node image (kindest/node:v1.34.0) 🖼
  ✓ Preparing nodes 📦
  ✓ Writing configuration 📜
  ✓ Starting control-plane 🕹️
  ✓ Installing CNI 🔌
  ✓ Installing StorageClass 💾
  Set kubectl context to "kind-budget-tracker-cluster"
  You can now use your cluster with:

  kubectl cluster-info --context kind-budget-tracker-cluster

  Not sure what to do next? 😅  Check out https://kind.sigs.k8s.io/docs/user/quick-start/
  ```

---

A local cluster using the current latest k8s' version was created and be accessed through it's local API server. You can run simply a `kubectl` command to validate clusters' health:


=== "Shell"
  ```shell
  kind cluster-info --context kind-budget-tracker-cluster
  ```
=== "Shell output"
  ```shell
  kubectl cluster-info
  Kubernetes control plane is running at https://127.0.0.1:50789
  CoreDNS is running at https://127.0.0.1:50789/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

  To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
  ```

!!! success "All set!"
    If you are able to reach your local k8s cluster, you can proceed by deployment the app to k8s.

### Build your local's container image

```shell
  docker build . -t budget-tracker-api:local
```

### Load your newest image to your local k8s cluster

```shell
  kind load docker-image budget-tracker-api:local --name budget-tracker-cluster
```

## Deploying the app through a helm chart

```shell
  make k8s-apply
```

<script src="https://giscus.app/client.js"
        data-repo="vsantos/budget-tracker-api-v2-discussions"
        data-repo-id="R_kgDOQApX1g"
        data-category="General"
        data-category-id="DIC_kwDOQApX1s4CwhAe"
        data-mapping="pathname"
        data-strict="0"
        data-reactions-enabled="1"
        data-emit-metadata="0"
        data-input-position="top"
        data-theme="catppuccin_frappe"
        data-lang="en"
        crossorigin="anonymous"
        async>
</script>