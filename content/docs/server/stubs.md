---
title: 'Stubs'
weight: 70
slug: stubs
aliases:
- /docs/stubs/
description: "Define precise responses for your gRPC APIs using FauxRPC stubs, enabling comprehensive testing and development."
icon: "design_services"
---

### Stub Definition

Stubs are defined using the `Stub` message in the `stubs.v1` package:

```protobuf
message Stub {
  StubRef ref = 1 [(buf.validate.field).required = true];
  oneof content {
    bytes proto = 2;
    string json = 3;
    Error error = 4;
  }
  // CEL rule to decide if this stub should be used for a given request.
  string active_if = 5;
  // Similar to the json attribute but is a CEL expression that returns the result.
  string cel_content = 6;
  int32 priority = 7 [(buf.validate.field).int32 = {
    gte: 0
    lte: 100
  }];
}
```

* **ref:**  A `StubRef` message that identifies the target method or type for this stub.
* **content:**  The response content, which can be:
    * **proto:** Raw protobuf bytes.
    * **json:** A JSON representation of the protobuf message.
    * **error:** An `Error` message to simulate an error response.
* **active_if:** A CEL expression that determines if this stub should be used based on the request.
* **cel_content:** A CEL expression that dynamically generates the response content.
* **priority:** An integer from 0 to 100 (higher is more preferred) to manage stub selection when multiple stubs match a request.


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

### Using CEL for Dynamic Responses

CEL (Common Expression Language) allows you to define dynamic stub behavior based on the request. You can use CEL in two ways:

* **active_if:**  Determine if a stub should be used based on a condition.
* **cel_content:**  Generate the response content dynamically.

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

### Conclusion

FauxRPC stubs provide a powerful mechanism for controlling API responses, enabling you to simulate various scenarios, test edge cases, and develop against your APIs without a real backend. By combining static responses, dynamic CEL expressions, and priority management, you can achieve fine-grained control over your mock services and streamline your development workflow.

The examples mentioned here are available in the [sudorandom/fauxrpc](https://github.com/sudorandom/fauxrpc/tree/main/example/stubs.petstore/) repo.
