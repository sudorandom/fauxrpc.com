---
title: Dashboard
---

Enhance your FauxRPC experience with the interactive dashboard, providing real-time insights into your server's operations.

To enable the dashboard, simply start FauxRPC with the `--dashboard` option:

```shell
fauxrpc run --schema=service.binpb --dashboard
```

Access the dashboard in your browser at `http://127.0.0.1:6660/fauxrpc`.

The dashboard provides:
*   📊 **Summary:** View overall server statistics.
    ![FauxRPC Dashboard](</dashboard.png>)
*   📜 **Request Log:** Live stream of all incoming requests.
    ![Request Log](</dashboard-request-log.png>)
*   📁 **Schema Browser:** Explore all Protobuf schemas loaded into the server. For more details on defining inputs, refer to the [Inputs page](../inputs/).
    ![Schema Browser 1](</dashboard-schema1.png>)
    ![Schema Browser 2](</dashboard-schema2.png>)
*   🔌 **Stubs:** Manage and view details of registered stubs. Learn more about stubs on the [Stubs page](../stubs/).
    ![Stubs](</dashboard-stubs.png>)
*   📚 **API Documentation:** Access auto-generated API documentation.
