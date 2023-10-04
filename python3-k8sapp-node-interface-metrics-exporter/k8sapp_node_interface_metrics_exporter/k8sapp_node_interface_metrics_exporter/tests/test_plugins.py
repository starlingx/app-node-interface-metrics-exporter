#
# Copyright (c) 2023 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#

from sysinv.tests.db import base as dbbase

from k8sapp_node_interface_metrics_exporter.common import constants as app_constants


class K8SAppNodeInterfaceMetricsExporterAppMixin(object):
    app_name = app_constants.HELM_APP_METRICS_EXPORTER
    path_name = app_name + '.tgz'

    # pylint: disable=invalid-name,useless-parent-delegation
    def setUp(self):
        super().setUp()

    def test_stub(self):
        # Replace this with a real unit test.
        pass


# Test Configuration:
# - Controller
# - IPv6
# - Ceph Storage
# - node-interface-metrics-exporter app
class K8SAppNodeInterfaceMetricsExporterControllerTestCase(
                                K8SAppNodeInterfaceMetricsExporterAppMixin,
                                dbbase.BaseIPv6Mixin,
                                dbbase.BaseCephStorageBackendMixin,
                                dbbase.ControllerHostTestCase):
    pass


# Test Configuration:
# - AIO
# - IPv4
# - Ceph Storage
# - node-interface-metrics-exporter app
class K8SAppNodeInterfaceMetricsExporterAIOTestCase(
                                K8SAppNodeInterfaceMetricsExporterAppMixin,
                                dbbase.BaseCephStorageBackendMixin,
                                dbbase.AIOSimplexHostTestCase):
    pass
