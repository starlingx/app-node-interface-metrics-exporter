#
# Copyright (c) 2023 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#
# All Rights Reserved.
#

""" System inventory App lifecycle operator."""

from oslo_log import log as logging

from sysinv.common import constants as sysinv_constants
from sysinv.helm import lifecycle_base as base
from sysinv.helm import lifecycle_utils

from k8sapp_node_interface_metrics_exporter.common import constants as app_constants

LOG = logging.getLogger(__name__)


