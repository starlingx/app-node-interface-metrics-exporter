#
# Copyright (c) 2023 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#

from oslo_log import log as logging

from k8sapp_node_interface_metrics_exporter.common import constants as app_constants

from sysinv.common import exception
from sysinv.helm import base


LOG = logging.getLogger(__name__)


class NodeInterfaceMetricsExporterHelm(base.BaseHelm):
    """Class to encapsulate helm operations for the metrics exporter chart."""

    SUPPORTED_NAMESPACES = base.BaseHelm.SUPPORTED_NAMESPACES + \
        [app_constants.HELM_NS_METRICS_EXPORTER]

    SUPPORTED_APP_NAMESPACES = {
        app_constants.HELM_APP_METRICS_EXPORTER: SUPPORTED_NAMESPACES
    }

    SERVICE_NAME = app_constants.HELM_APP_METRICS_EXPORTER

    SUPPORTED_COMPONENT_OVERRIDES = [
        app_constants.HELM_COMPONENT_LABEL_VALUE_PLATFORM,
        app_constants.HELM_COMPONENT_LABEL_VALUE_APPLICATION
    ]

    CHART = app_constants.HELM_CHART_METRICS_EXPORTER

    def get_namespaces(self):
        return self.SUPPORTED_NAMESPACES

    def get_overrides(self, namespace=None):
        overrides = {
            app_constants.HELM_NS_METRICS_EXPORTER: {}
        }

        if namespace in self.SUPPORTED_NAMESPACES:
            return overrides[namespace]

        if namespace:
            raise exception.InvalidHelmNamespace(chart=self.CHART,
                                                 namespace=namespace)
        return overrides
