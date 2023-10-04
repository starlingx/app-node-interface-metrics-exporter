#
# Copyright (c) 2023 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#

# Namespace to deploy the application
HELM_NS_METRICS_EXPORTER = 'node-interface-metrics-exporter'

# Application Name
HELM_APP_METRICS_EXPORTER = 'node-interface-metrics-exporter'

# Chart Name
HELM_CHART_METRICS_EXPORTER = 'node-interface-metrics-exporter'

# Application component label
HELM_LABEL_PARAMETER = 'podLabels'
HELM_COMPONENT_LABEL = 'app.starlingx.io/component'
HELM_COMPONENT_LABEL_VALUE_PLATFORM = 'platform'
HELM_COMPONENT_LABEL_VALUE_APPLICATION = 'application'
