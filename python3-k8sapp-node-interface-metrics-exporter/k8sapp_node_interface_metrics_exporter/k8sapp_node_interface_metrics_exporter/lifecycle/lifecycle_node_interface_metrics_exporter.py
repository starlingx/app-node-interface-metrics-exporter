#
# Copyright (c) 2023-2024 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#
# All Rights Reserved.
#

""" System inventory App lifecycle operator."""

from oslo_log import log as logging

import yaml

from sysinv.common import constants as inv_constants
from sysinv.common import exception
from sysinv.common import kubernetes
from sysinv.common import utils as cutils
from sysinv.helm import lifecycle_base as base
from sysinv.helm.lifecycle_constants import LifecycleConstants

from k8sapp_node_interface_metrics_exporter.common import constants as app_constants

LOG = logging.getLogger(__name__)


class NodeInterfaceMetricsExporterAppLifecycleOperator(base.AppLifecycleOperator):
    def app_lifecycle_actions(
        self, context, conductor_obj, app_op, app, hook_info
    ):
        """Perform lifecycle actions for an operation

        :param context: request context, can be None
        :param conductor_obj: conductor object, can be None
        :param app_op: AppOperator object
        :param app: AppOperator.Application object
        :param hook_info: LifecycleHookInfo object

        """
        if hook_info.lifecycle_type == LifecycleConstants.APP_LIFECYCLE_TYPE_FLUXCD_REQUEST:
            if hook_info.operation == inv_constants.APP_APPLY_OP:
                if hook_info.relative_timing == LifecycleConstants.APP_LIFECYCLE_TIMING_POST:
                    return self.post_apply(app_op, app, hook_info)

                if hook_info.relative_timing == LifecycleConstants.APP_LIFECYCLE_TIMING_PRE:
                    # on pre apply hook adding Label
                    self.assign_host_label(app_op)

        if hook_info.lifecycle_type == LifecycleConstants.APP_LIFECYCLE_TYPE_OPERATION:
            if hook_info.operation == inv_constants.APP_REMOVE_OP:
                if hook_info.relative_timing == LifecycleConstants.APP_LIFECYCLE_TIMING_POST:
                    # on post remove hook removing labels
                    self.remove_host_labels(app_op)
                    return self.post_remove(app)

        super(
            NodeInterfaceMetricsExporterAppLifecycleOperator, self
        ).app_lifecycle_actions(context, conductor_obj, app_op, app, hook_info)

    def post_apply(self, app_op, app, hook_info):
        if LifecycleConstants.EXTRA not in hook_info:
            raise exception.LifecycleMissingInfo("Missing {}".format(LifecycleConstants.EXTRA))
        if LifecycleConstants.RETURN_CODE not in hook_info[LifecycleConstants.EXTRA]:
            raise exception.LifecycleMissingInfo(
                "Missing {} {}".format(LifecycleConstants.EXTRA, LifecycleConstants.RETURN_CODE))

        # raise a specific exception to be caught by the
        # retry decorator and attempt a re-apply
        if not hook_info[LifecycleConstants.EXTRA][LifecycleConstants.RETURN_CODE] and \
                not app_op.is_app_aborted(app.name):
            LOG.info("%s app failed applying. Retrying." % str(app.name))
            raise exception.ApplicationApplyFailure(name=app.name)

        dbapi_instance = app_op._dbapi
        db_app_id = dbapi_instance.kube_app_get(app.name).id

        client_core = app_op._kube._get_kubernetesclient_core()
        component_constant = (
            app_constants.HELM_COMPONENT_LABEL
        )

        # chart overrides
        chart_overrides = self._get_helm_user_overrides(
            dbapi_instance, db_app_id
        )

        override_label = {}

        # Namespaces variables
        namespace = client_core.read_namespace(
            app_constants.HELM_NS_METRICS_EXPORTER
        )

        # Old namespace variable
        old_namespace_label = (
            namespace.metadata.labels.get(component_constant)
            if component_constant in namespace.metadata.labels
            else None
        )

        if component_constant in chart_overrides:
            # User Override variables
            dict_chart_overrides = yaml.safe_load(chart_overrides)
            override_label = dict_chart_overrides.get(component_constant)

        if override_label == "application":
            namespace.metadata.labels.update(
                {component_constant: "application"}
            )
            app_op._kube.kube_patch_namespace(
                app_constants.HELM_NS_METRICS_EXPORTER, namespace
            )
        elif override_label == "platform":
            namespace.metadata.labels.update({component_constant: "platform"})
            app_op._kube.kube_patch_namespace(
                app_constants.HELM_NS_METRICS_EXPORTER, namespace
            )
        elif not override_label:
            namespace.metadata.labels.update({component_constant: "platform"})
            app_op._kube.kube_patch_namespace(
                app_constants.HELM_NS_METRICS_EXPORTER, namespace
            )
        else:
            LOG.info(
                f"WARNING: Namespace label {override_label} not supported"
            )

        namespace_label = namespace.metadata.labels.get(component_constant)
        if old_namespace_label != namespace_label:
            self._delete_ineterface_metrics_exporter_pods(app_op, client_core)

    def post_remove(self, app):
        LOG.debug(
            "Executing post_remove for {} app".format(
                app_constants.HELM_APP_METRICS_EXPORTER
            )
        )
        cmd = [
            "kubectl",
            "--kubeconfig",
            kubernetes.KUBERNETES_ADMIN_CONF,
            "delete",
            "namespace",
            app_constants.HELM_NS_METRICS_EXPORTER,
        ]
        stdout, stderr = cutils.trycmd(*cmd)
        LOG.debug(
            "{} app: cmd={} stdout={} stderr={}".format(
                app.name, cmd, stdout, stderr
            )
        )

    def _get_helm_user_overrides(self, dbapi_instance, db_app_id):
        try:
            overrides = dbapi_instance.helm_override_get(
                app_id=db_app_id,
                name=app_constants.HELM_APP_METRICS_EXPORTER,
                namespace=app_constants.HELM_NS_METRICS_EXPORTER,
            )
        except exception.HelmOverrideNotFound:
            values = {
                "name": app_constants.HELM_APP_METRICS_EXPORTER,
                "namespace": app_constants.HELM_NS_METRICS_EXPORTER,
                "db_app_id": db_app_id,
            }
            overrides = dbapi_instance.helm_override_create(values=values)
        return overrides.user_overrides or ""

    def _delete_ineterface_metrics_exporter_pods(self, app_op, client_core):
        # pod list
        system_pods = client_core.list_namespaced_pod(
            app_constants.HELM_NS_METRICS_EXPORTER
        )

        # On namespace label change delete pods to force restart
        for pod in system_pods.items:
            app_op._kube.kube_delete_pod(
                name=pod.metadata.name,
                namespace=app_constants.HELM_NS_METRICS_EXPORTER,
                grace_periods_seconds=0,
            )

    def assign_host_label(self, app_op):
        """
        function to assign labels
        """
        hosts = app_op._dbapi.ihost_get_list()
        label_key, label_value = app_constants.NODE_LABEL.split('=')
        label_dict = {"label_key": label_key, "label_value": label_value}
        for host in hosts:
            # subfunctions can have values like "controller,worker", "worker"
            # "controller", "storage"
            # checking if contains worker "worker" in "controller,worker"
            if inv_constants.WORKER in host.subfunctions:
                # assign Label
                LOG.info("assign label Node={} has role={}".format(host.hostname, host.subfunctions))
                try:
                    app_op._dbapi.label_create(
                        host.id, {"host_id": host.id, **label_dict}
                    )
                except exception.HostLabelAlreadyExists:
                    pass
                app_op._update_kubernetes_labels(host.hostname, {label_key: label_value})

    def remove_host_labels(self, app_op):
        """
        function to remove labels
        """
        hosts = app_op._dbapi.ihost_get_list()
        for host in hosts:
            # subfunctions can have values like "controller,worker", "worker"
            # "controller", "storage"
            # checking if contains worker "worker" in "controller,worker"
            if inv_constants.WORKER in host.subfunctions:
                LOG.info("remove label Node={} has role={}".format(host.hostname, host.subfunctions))
                # remove Label
                lbl_obj = app_op._find_label(host.uuid, app_constants.NODE_LABEL)
                if lbl_obj:
                    app_op._dbapi.label_destroy(lbl_obj.uuid)
                    app_op._update_kubernetes_labels(host.hostname, {lbl_obj.label_key: None})
