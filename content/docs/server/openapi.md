---
title: 'OpenAPI Support'
weight: 35
slug: openapi
description: "Serve and stub OpenAPI operations with generated, schema-aware HTTP responses."
icon: "api"
---

FauxRPC can create a working HTTP mock server directly from an OpenAPI 3 specification. OpenAPI and Protobuf schemas use the same `--schema` option and can be served by the same FauxRPC process.

## Load an OpenAPI schema

Pass a YAML or JSON document to `fauxrpc run`:

```shell
fauxrpc run --schema=./openapi.yaml
```

The schema can be a local file, an HTTP(S) URL, or a directory containing OpenAPI documents. You can repeat `--schema` to combine sources:

```shell
fauxrpc run \
  --schema=./service.binpb \
  --schema=./openapi.yaml
```

FauxRPC honors the server base path declared by the OpenAPI document. For example, a schema with `servers: [{ url: /api/v3 }]` serves `/pets` at `/api/v3/pets`.

## Generated responses

When an operation has no matching stub, FauxRPC selects a declared response and generates data from it. It prefers a `200` response, then `201`, then a default or another declared response.

Generation supports:

* Objects, arrays, and primitive schemas
* Examples and defaults
* Enums, numeric and string constraints, and common formats
* `allOf`, `oneOf`, and `anyOf`
* Recursive schemas, bounded by `--depth`
* Response headers declared in the selected response

Explicit examples and defaults are returned as written. Other faux values vary between requests by default.

### Reproducible generated data

Use `--static-seed` when tests or snapshots need repeatable generated values:

```shell
fauxrpc run --schema=./openapi.yaml --static-seed
```

With this option, each OpenAPI operation and Protobuf RPC method receives a stable, identity-derived seed. Repeated unstubbed calls therefore produce the same generated values. Stubs and explicit schema examples or defaults are unaffected.

Without `--static-seed`, unstubbed OpenAPI and Protobuf responses receive fresh generated values on each request.

### Disable generated fallbacks

Use `--only-stubs` when every OpenAPI response must come from a configured stub:

```shell
fauxrpc run \
  --schema=./openapi.yaml \
  --stubs=./openapi-stubs.yaml \
  --only-stubs
```

Requests are still matched and validated against the OpenAPI operation. If no stub matches the valid request, FauxRPC returns HTTP `501 Not Implemented` with a `no matching stub` message.

## Request validation

FauxRPC validates incoming path parameters, query parameters, headers, and request bodies against the matched operation before selecting a stub or generating a response. Invalid requests receive an HTTP `400` response with the validation error.

FauxRPC does not currently implement authentication callbacks for OpenAPI security schemes. For a self-contained local mock, override secured operations with `security: []` in the mock copy of the specification.

## Interactive documentation

Every loaded OpenAPI document is available through an interactive Scalar API reference:

```text
http://127.0.0.1:6660/fauxrpc/openapi-docs/
```

The document served to Scalar points its primary server at the running FauxRPC instance and preserves the schema's base path. Disable documentation with `--no-doc-page`.

## OpenAPI stubs

Use `--stubs` to return precise responses for selected operations. A stub can target an `operationId`:

```yaml
stubs:
  - name: The answer to pets, life, and everything
    target:
      operationId: getPetById
    match:
      pathParams:
        petId: "42"
    response:
      status: 200
      headers:
        X-Pet-Mood: existential
      body:
        id: 42
        name: Deep Thought
        status: available
```

It can also target a path and HTTP method:

```yaml
stubs:
  - name: Missing pet
    target:
      path: /pet/{petId}
      httpMethod: GET
    match:
      pathParams:
        petId: "404"
    response:
      status: 404
      body:
        message: Pet not found
```

Matches can use:

* `pathParams`
* `queryParams`
* `headers`
* `bodyExpression`, a [GJSON](https://github.com/tidwall/gjson) path whose selected value must be truthy

When several stubs match, FauxRPC prefers the highest priority, then the most specific match, then insertion order.

Load the schema and stubs together:

```shell
fauxrpc run \
  --schema=./openapi.yaml \
  --stubs=./openapi-stubs.yaml
```

See [Stubs](/docs/stubs/) for Protobuf stub features and additional configuration examples.
