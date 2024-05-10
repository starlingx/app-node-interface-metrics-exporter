# Node Interface Metricse Exporter Helm Chart

This Helm Chart comes with everything that is needed to run Node Interface
Metrics Exporter app.

## How to install

Create the namespace where Helm should install the components with

```console
$ kubectl create namespace "node-interface-metrics-exporter"
```

This chart can be installed using the default values with:

```console
$ helm upgrade --install node-interface-metrics-exporter stx-platform/node-interface-metrics-exporter
```

It is possible to replace the default "args" values of values.yaml file to
change or add new args to Node interface metrics exporter.

You may want to override values.yaml file use the following command:

```console
$ helm upgrade --install node-interface-metrics-exporter -f values.yaml stx-platform/node-interface-metrics-exporter
```

### Delete Chart

If you want to delete your Chart, use this command

```console
$ helm uninstall node-interface-metrics-exporter
```

If you want to delete the namespace, use this command

```console
$ kubectl delete namespace node-interface-metrics-exporter
```

For more information about installing and using Helm, see the
[Helm Docs](https://helm.sh/docs/). For a quick introduction to Charts, 
see the [Chart Guide](https://helm.sh/docs/topics/charts/).
