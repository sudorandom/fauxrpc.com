---
title: Dashboard
---

Inspect live multi-protocol traffic, registered stubs, and Protobuf schemas from the FauxRPC web dashboard.

Pass `--dashboard` to `fauxrpc run`:

```shell
fauxrpc run --schema=service.binpb --dashboard
```

Open `http://127.0.0.1:6660/fauxrpc` in your browser.

### Live Telemetry & Overview

The top application bar tracks real-time throughput, total requests, error rates, loaded services, active stubs, and server uptime. Expand the status popover to view runtime diagnostics or reset request counters.

{{< img src="/images/dashboard.png" alt="FauxRPC Dashboard Overview and Telemetry" width="800" />}}

### Request Log & Inspector

Incoming calls stream into the Request Log over Server-Sent Events as they arrive. Filter rows by service, method, payload text, status outcome (`OK` or `Errors`), or client protocol (`gRPC`, `gRPC-Web`, `Connect`, `HTTP`). Selecting any row opens the right-hand inspector with formatted request and response payloads, headers, matched stubs, and a generated YAML stub ready to copy, download, or register directly into the running server.

{{< img src="/images/dashboard-request-log.png" alt="Request Log and Payload Inspector" width="800" />}}

### Schema Browser

Navigate every registered Protobuf package and file in a split tree-and-source view. Each file entry displays service and message counts alongside syntax-highlighted `.proto` source reconstructed from the active descriptors. See the [Inputs page](../inputs/) for supported schema sources.

{{< img src="/images/dashboard-schema1.png" alt="Schema Browser - Protobuf Service Definition" width="800" />}}
{{< img src="/images/dashboard-schema2.png" alt="Schema Browser - REST Annotations" width="800" />}}

### Stubs Manager

View, filter, create, and edit registered mock stubs in real time. The Stubs view displays each stub's mode (`CEL`, `JSON`, `ERROR`, `STREAM`), target service or method, `active_if` CEL rule, priority, and equivalent YAML configuration. Read the [Stubs page](../stubs/) for details on authoring rules.

{{< img src="/images/dashboard-stubs.png" alt="Stubs Inspector and Editor" width="800" />}}

### API Documentation

FauxRPC generates interactive Protobuf service documentation (`/fauxrpc/docs`) and OpenAPI documentation (`/fauxrpc/openapi.html`) directly from loaded schemas. Both views are accessible from the top navigation bar.

