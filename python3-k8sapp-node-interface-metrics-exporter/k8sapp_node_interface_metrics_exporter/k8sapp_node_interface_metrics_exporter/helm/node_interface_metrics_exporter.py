#
# Copyright (c) 2023 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#

from k8sapp_node_interface_metrics_exporter.common import constants as app_constants

from oslo_log import log as logging
from sysinv.common import constants
from sysinv.common import exception
from sysinv.common import kubernetes
from sysinv.db import api as sys_dbapi
from sysinv.helm import base


LOG = logging.getLogger(__name__)

class NodeInterfaceMetricsExporterHelm(base.FluxCDBaseHelm):
    """Class to encapsulate helm operations for the metrics exporter chart."""

    SUPPORTED_NAMESPACES = base.FluxCDBaseHelm.SUPPORTED_NAMESPACES + \
        [app_constants.HELM_NS_PM]

    SUPPORTED_APP_NAMESPACES = {
        app_constants.HELM_APP_PM: SUPPORTED_NAMESPACES
    }

    SERVICE_NAME = app_constants.HELM_APP_PM

    SUPPORTED_COMPONENT_OVERRIDES = [
        app_constants.HELM_COMPONENT_LABEL_VALUE_PLATFORM,
        app_constants.HELM_COMPONENT_LABEL_VALUE_APPLICATION
    ]

    DEFAULT_AFFINITY = app_constants.HELM_COMPONENT_LABEL_VALUE_PLATFORM

    def get_namespaces(self):
        return self.SUPPORTED_NAMESPACES

    def get_overrides(self, namespace=None):
        dbapi_instance = dbapi.get_instance()
        db_app = dbapi_instance.kube_app_get(app_constants.HELM_APP_PM)

        # User chart overrides
        chart_overrides = self._get_user_helm_overrides(
            dbapi_instance,
            db_app,
            app_constants.HELM_NS_PM,
            'user_overrides')

        user_affinity = chart_overrides.get(app_constants.HELM_COMPONENT_LABEL,
                                            self.DEFAULT_AFFINITY)

        if user_affinity in self.SUPPORTED_COMPONENT_OVERRIDES:
            affinity = user_affinity
        else:
            LOG.warn(f"User override value {user_affinity} "
                     f"for {app_constants.HELM_COMPONENT_LABEL} is invalid, "
                     f"using default value {self.DEFAULT_AFFINITY}")
            affinity = self.DEFAULT_AFFINITY

        overrides = {
            app_constants.HELM_NS_PM: {
                app_constants.HELM_LABEL_PARAMETER: {
                    app_constants.HELM_COMPONENT_LABEL: affinity
                }
            }
        }

        if namespace in self.SUPPORTED_NAMESPACES:
            return overrides[namespace]

        if namespace:
            raise exception.InvalidHelmNamespace(chart=self.CHART,
                                                 namespace=namespace)
        return overrides

    @staticmethod
    def _get_user_helm_overrides(dbapi_instance, app, chart, namespace,
                            type_of_overrides):
        """Helper function for querying helm overrides from db."""
        helm_overrides = {}
        try:
            overrides = dbapi_instance.helm_override_get(
                app_id=app.id,
                name=chart,
                namespace=namespace,
            )[type_of_overrides]

            if isinstance(overrides, str):
                helm_overrides = yaml.safe_load(overrides)
        except exception.HelmOverrideNotFound:
            LOG.debug("Overrides for this chart not found, nothing to be done.")
        return helm_overrides

