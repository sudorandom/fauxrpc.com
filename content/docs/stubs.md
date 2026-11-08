---
title: 'Stubs'
weight: 51
slug: stubs
aliases:
- /docs/server/stubs/
- /docs/stubs/
description: "Define precise responses for Protobuf RPCs and OpenAPI operations using FauxRPC stubs."
icon: "design_services"
---

FauxRPC stubs allow you to define precise responses for Protobuf RPC calls and OpenAPI operations. Stubs are supported by:

*   **`fauxrpc run`**: The server will respond with the stubbed data when a matching request is received.
*   **`fauxrpc curl`**: The client will use the stubbed data as the request payload.
*   **`fauxrpc generate`**: The generator will output the stubbed data to stdout.

OpenAPI stubs use operation, path, and HTTP request matching instead of the Protobuf fields described below. See [OpenAPI Stubs](/docs/server/openapi/#openapi-stubs) for their format and examples.

### Stub Configuration

Stubs can be defined using YAML or JSON. The configuration object has the following properties:

* **target:** The fully qualified name of the service method or message type (e.g., `connectrpc.eliza.v1.ElizaService/Say`).
* **id:** A unique identifier for the stub.
* **priority:** An integer from 0 to 100 (higher is more preferred).
* **active_if:** A CEL expression that determines if this stub should be used.
* **content:** A JSON object representing the response message.
* **cel_content:** A CEL expression that dynamically generates the response content.
* **stream:** Configuration for streaming responses.
* **error_code:** gRPC status code for error responses.
* **error_message:** Error message for error responses.

A stub must define exactly one of `content`, `cel_content`, `stream`, or `error_code`/`error_message`.


### Adding Stubs

While that protobuf message is how stubs are defined, you can add stubs using the `fauxrpc stub add` CLI command:

```shell
fauxrpc stub add <target> [flags]
```

where `<target>` is the fully qualified protobuf method or type.

**Flags:**

* `--id`:  A unique identifier for the stub.
* `--json`:  JSON response content.
* `--error-message`:  Error message for error responses.
* `--error-code`:  gRPC error code for error responses.
* `--active-if`:  CEL expression for conditional activation.
* `--priority`:  Stub priority.
* `--cel`:  CEL expression for dynamic content generation.

**Examples:**

* **Return a static JSON response:**

```shell
fauxrpc stub add connectrpc.eliza.v1.ElizaService/Say --json '{"sentence": "Test from the CLI"}'
```

* **Return an error:**

```shell
fauxrpc stub add buf.registry.owner.v1.UserService/CreateUsers --error-message="hello!" --error-code=5 
```

### Defining Stubs in YAML

You can also define stubs in YAML (or JSON) files. This allows for more complex scenarios and better organization.

**Example (dynamic content with CEL):**

```yaml
---
stubs:
- id: add-pet
  target: io.swagger.petstore.v2.PetService/AddPet
  cel_content: req 
```

This stub will return the request itself as the response.

**Example (static response):**

```yaml
---
stubs:
- id: find-pets-by-status
  target: io.swagger.petstore.v2.PetService/FindPetsByStatus
  content:
    pets:
    - id: 1
      category:
        id: 1
        name: cat
      name: Whiskers
      photo_urls:
      - https://cataas.com/cat
      tags:
      - id: 1
        name: cute
      - id: 2
        name: kid-friendly
      status: available
```

This stub will always return a predefined list of pets.

**Example (combining `active_if`, `priority`, and `cel_content`):**

```yaml
---
stubs:
- id: get-pets-by-id-id-1
  target: io.swagger.petstore.v2.PetService/GetPetByID
  active_if: req.pet_id == 1
  priority: 100
  content:
    id: 1
    category:
      id: 1
      name: cat
    name: Whiskers
    photo_urls:
    - https://cataas.com/cat
    tags:
    - id: 1
      name: cute
    - id: 2
      name: kid-friendly
    status: available
- id: get-pets-by-id-default
  target: io.swagger.petstore.v2.PetService/GetPetByID
  cel_content: |
    {
        'id': req.pet_id,
        'category': {'id': gen, 'name': 'gen'},
        'name': gen,
        'photo_urls': [gen, gen],
        'tags': [{'id': gen, 'name': gen}],
        'status': gen
    }
```

In this example:

* The first stub provides a specific response for when `pet_id` is 1.
* The second stub uses `cel_content` and the `gen` value to generate a dynamic response for any other `pet_id`.
* This demonstrates how to handle specific cases with high-priority stubs while providing a general fallback with dynamic content.

### Streaming Stubs

FauxRPC supports streaming responses, allowing you to simulate server-side streaming or bidirectional streaming scenarios.

**Stream Configuration:**

* **items:** A list of items to stream.
* **repeated:** If `true`, the sequence of items is repeated indefinitely (or until `done_after` is reached).
* **done_after:** The total duration for the stream to run (e.g., `30s`, `1m`).

**Stream Item Configuration:**

Each item in the stream can have:

* **content:** A JSON object representing the response message.
* **cel_content:** A CEL expression that evaluates to the response message.
* **error:** An object with `code` (integer) and `message` (string) to return an error.
* **delay:** Duration to wait before sending this message (e.g., `100ms`, `1s`).

**Example (basic stream):**

```yaml
stubs:
  - id: introduce-basic
    target: connectrpc.eliza.v1.ElizaService/Introduce
    stream:
      items:
        - content: { sentence: "Hello, I am Eliza." }
          delay: 100ms
        - content: { sentence: "I am here to listen." }
          delay: 500ms
        - content: { sentence: "What is on your mind?" }
          delay: 500ms
```

**Example (repeated stream with active_if):**

```yaml
stubs:
  - id: introduce-repeated
    target: connectrpc.eliza.v1.ElizaService/Introduce
    active_if: req.name == "repeat"
    priority: 10
    stream:
      repeated: true
      done_after: 30s
      items:
        - content: { sentence: "This message repeats." }
          delay: 1s
```

**Example (stream with dynamic content using CEL):**

```yaml
stubs:
  - id: introduce-cel
    target: connectrpc.eliza.v1.ElizaService/Introduce
    active_if: req.name == "cel"
    priority: 10
    stream:
      items:
        - cel_content: "{'sentence': 'Hello ' + req.name}"
          delay: 200ms
        - cel_content: "{'sentence': 'Nice to meet you, ' + req.name}"
          delay: 200ms
```

**Example (stream ending with an error):**

```yaml
stubs:
  - id: introduce-error
    target: connectrpc.eliza.v1.ElizaService/Introduce
    active_if: req.name == "error"
    priority: 10
    stream:
      items:
        - content: { sentence: "Something is about to go wrong..." }
          delay: 200ms
        - error:
            code: 13 # INTERNAL
            message: "Something went wrong"
```

### Using CEL for Dynamic Responses

CEL (Common Expression Language) allows you to define dynamic stub behavior based on the request. You can use CEL in two ways:

* **active_if:**  Determine if a stub should be used based on a condition.
* **cel_content:**  Generate the response content dynamically.

#### Available Variables

The following variables are available in the CEL environment:

*   **`req`**: The request message. You can access fields using dot notation (e.g., `req.name`).
*   **`now`**: The current timestamp (google.protobuf.Timestamp).
*   **`gen`**: A helper for generating random data (used in `cel_content` primarily).

**Example (using `active_if`):**

```yaml
stubs:
- id: add-pet-active-if-example
  target: io.swagger.petstore.v2.PetService/AddPet
  json: '{"id": 1}'
  active_if: 'req.name == "Fluffy"' 
```

This stub will only be active if the `name` field in the request is "Fluffy".

**Example (using `cel_content`):**

```yaml
stubs:
- id: add-pet-cel-content-example
  target: io.swagger.petstore.v2.PetService/AddPet
  cel_content: '{"id": req.id, "name": req.name + " the Great"}' 
```

This stub will generate a response where the `id` is taken from the request and the `name` is the request's `name` with " the Great" appended.

### Stub Priority

When multiple stubs match a request, FauxRPC uses the `priority` field to determine which stub to use. Higher priority stubs (100 being the highest) are preferred over lower priority stubs.

This allows you to define default responses with lower priority and more specific responses with higher priority.

### Recording Stubs from Upstream Traffic

Instead of writing YAML/JSON stub files by hand or adding them one-by-one via the CLI, you can automatically record real traffic passing through FauxRPC when it is run in proxy mode.

By starting the server with `--proxy-to` and `--record-dir`, gRPC/Connect requests and responses are captured and saved to disk in the structured stub format. See the [Proxy & Record](/docs/server/proxy-and-record/) guide for detailed information.

### Replaying Recorded Stubs

Once stubs are recorded in a directory, you can load them back into FauxRPC offline:

```shell
fauxrpc run \
  --schema=buf.build/connectrpc/eliza \
  --stubs=stubs/
```

### Conclusion

FauxRPC stubs provide a powerful mechanism for controlling API responses, enabling you to simulate various scenarios, test edge cases, and develop against your APIs without a real backend. By combining static responses, dynamic CEL expressions, and priority management, you can achieve fine-grained control over your mock services and streamline your development workflow.

The examples mentioned here are available in the [sudorandom/fauxrpc](https://github.com/sudorandom/fauxrpc/tree/main/example/stubs.petstore/) repo.
