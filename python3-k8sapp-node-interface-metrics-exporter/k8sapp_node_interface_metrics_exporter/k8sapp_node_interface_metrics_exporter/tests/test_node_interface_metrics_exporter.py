#
# Copyright (c) 2023 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#

from sysinv.db import api as dbapi
from sysinv.tests.db import base as dbbase
from sysinv.tests.db import utils as dbutils
from sysinv.tests.helm import base

from k8sapp_node_interface_metrics_exporter.tests import test_plugins


class NodeInterfaceMetricsExporterTestCase(
                    test_plugins.K8SAppNodeInterfaceMetricsExporterAppMixin,
                    base.HelmTestCaseMixin):
    """Test Class for Node interface metris exporter app."""

    def setUp(self):
        super().setUp()
        self.app = dbutils.create_test_app(
                    name='node-interface-metrics-exporter')
        self.dbapi = dbapi.get_instance()


class NodeInterfaceMetricsExporterTestCaseDummy(
                            NodeInterfaceMetricsExporterTestCase,
                            dbbase.ProvisionedControllerHostTestCase):
    """Dummy Class to pass the zuul."""

    def test_dummy(self):
        pass
