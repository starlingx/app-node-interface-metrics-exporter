=======================================
Node Interface Metrics Exporter API v1
=======================================

This Section contains the list of APIs of Node Interface Metrics Exporter.

The Application Pod retrieves the node interface information and statistics by making a call to the following REST APIs.

The server address is the node’s host name or IP address, and it uses port 9110.

The output is presented in the Open Metrics format, providing information and statistics for devices that are in the UP state.


---------------------
Open Metrics API List
---------------------

*****************************************
Get statistics of all devices in the node
*****************************************

.. rest_method:: GET /metrics

-   ``/metrics`` provides the information and statistics of all Physical
    Function devices.

-   It provides the information and statistics of Virtual Function devices
    that belong to the Physical Function devices.

-   It also provides the POD information if the Virtual Function device
    is used by any POD (matched by PciAddr of Virtual Function device).

**Normal response codes**

200

**Error response codes**

computeFault (400, 500, …), serviceUnavailable (503), badRequest (400), badMethod (405), itemNotFound (404)

**Request parameters**

None

**Response parameters**

Node statistics in Open Metrics format.

.. include:: samples/openmet_complete_pf_vfs.txt
    :literal:


****************************************************
Get Statistics of a particular device by device_name
****************************************************

.. rest_method:: GET /metrics/device/{device_name}

Get the information and statistics of a particular device in the Node
identified by ``{device_name}``.

-   If the input device name is Physical Function,
    ``/metrics/device/{device_name}`` provides its device information and
    statistics. It provides device information and statistics of all
    the Virtual Functions that belong to the input device. It also provides
    the POD information if Virtual Function device is used by any
    POD (matched by ``pci-address`` of Virtual Function device).

-   Input Virtual Function device name is not supported as multiple Virtual
    Functions have the same device name in different pods.

**Normal response codes**

200

**Error response codes**

computeFault (400, 500, …), serviceUnavailable (503), badRequest (400), badMethod (405), itemNotFound (404)

**Request parameters**

.. csv-table::
    :header: "Parameter", "Style", "Type", "Description"
    :widths: 20, 20, 20, 60

    "device_name","URI","string","The unique name of Physical Function device"

**Response parameters**

Device statistics in Open Metrics format.

.. include:: samples/openmet_complete_pf_vfs.txt
    :literal:

****************************************************
Get Statistics of a particular device by pci_address
****************************************************

.. rest_method:: GET /metrics/pci-addr/{pci_address}

Get statistics of the particular device in the node identified
by ``{pci_address}``.

-   If the input ``pci-address`` belongs to the Physical Function device,
    ``/metrics/pci-addr/{pci_address}`` provides its device information and
    statistics. It also provides statistics of all the Virtual
    Function devices that belong to the Physical Function device.

-   If the input ``pci-address`` belongs to the Virtual Function device, it
    provides its device information and statistics. It provides the POD
    information as well if the Virtual Function device is used by any POD
    (matched by ``pci-address`` of the Virtual Function device).

**Normal response codes**

200

**Error response codes**

computeFault (400, 500, …), serviceUnavailable (503), badRequest (400), badMethod (405), itemNotFound (404)

**Request parameters**

.. csv-table::
    :header: "Parameter", "Style", "Type", "Description"
    :widths: 20, 20, 20, 80

    "pci_address","URI","csapi:UUID","Unique identifier of the pci-address of the Physical Function or the Virtual Function devices"

**Response parameters**

Device statistics in Open Metrics format.

- When input is ``pci-address`` of Physical Function device:

.. include:: samples/openmet_complete_pf_vfs.txt
    :literal:

- When input is ``pci-address`` of the Virtual Function device:

.. include:: samples/openmet_pciaddr_vf.txt
    :literal:
