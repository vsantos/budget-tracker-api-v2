# Helm chart

This app can be installed as a "helm chart" for Kubernetes' environments.

## Generating manifests

 In order to do it so, you can generate the manifests locally through:

```shell
helm template --release-name local-dev ./helm/budget-tracker > manifests.yaml

cat manifests.yaml
```

## Deploying manifests

After inspecting and validating those manifests, they can be applied through:

```shell
kubectl apply -f manifests.yaml
```

### Testing manifests

You can run `helm-unittest` to validate budget-tracker's chart:

```shell
make helm-test
```

{% include "discussions.md" %}
