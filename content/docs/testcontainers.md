---
title: 'Testcontainers'
weight: 50
slug: testcontainers
description: "Simplify your gRPC testing by using FauxRPC with Testcontainers for lightweight, isolated environments."
icon: "deployed_code"
---

Testing gRPC services can be tricky. You often need a real server running, which can introduce complexity and slow down your tests. Using FauxRPC with [Testcontainers](https://testcontainers.com/) can simplify gRPC mocking.

To address challenges with testing while using gRPC services, we can leverage the power of [Testcontainers](https://testcontainers.com/), a library that lets you run throwaway, lightweight instances of common databases, web browsers, or any other application that can run in a Docker container. This allows you to easily integrate these dependencies into your automated tests, providing a consistent and reliable testing environment. By using Testcontainers, you can ensure that your tests are always running against a known and controlled version of your dependencies, avoiding inconsistencies and unexpected behavior while also simplifying test setup/teardown.

While Testcontainers provides the infrastructure, [FauxRPC](https://fauxrpc.com) takes care of the mocking itself. FauxRPC is a tool that generates fake gRPC, gRPC-Web, Connect, and REST servers from your Protobuf definitions. By combining it with Testcontainers, you gain a lightweight, isolated environment for testing your gRPC clients without relying on a real server implementation. I've made a package to make this simpler using Go but the same could be done for other languages that Testcontainers supports.

## Show by Example

### 1. Setting up the Container
```go
container, err := fauxrpctestcontainers.Run(ctx, "docker.io/sudorandom/fauxrpc:latest")
// ... error handling ...
t.Cleanup(func() { container.Terminate(context.Background()) })
```

This snippet starts a FauxRPC container using the `fauxrpctestcontainers.Run` function. The `t.Cleanup` function ensures the container is terminated after the test, keeping your testing environment clean.

### 2. Registering the Protobuf Definition
```go
container.MustAddFileDescriptor(ctx, elizav1.File_connectrpc_eliza_v1_eliza_proto)
```

You register your Protobuf file descriptor with the container. This lets FauxRPC understand the structure of your gRPC service. Now you have a fully functional FauxRPC service that mimics the services in the file descriptor that you gave it. The data is all randomly generated. Now let's test it a bit.

### 3. Making gRPC Calls
```go
baseURL := container.MustBaseURL(ctx)
elizaClient := elizav1connect.NewElizaServiceClient(http.DefaultClient, baseURL)
resp, err := elizaClient.Say(ctx, connect.NewRequest(&elizav1.SayRequest{
    Sentence: "testing!",
}))
// ... error handling and assertions ...
```

This code is getting base URL of the FauxRPC server running in the container and creating a gRPC client (using ConnectRPC). ConnectRPC isn't a requirement. You can use grpc-go instead. With this client, you can make calls to your gRPC service as you would in a real environment. In this setup, FauxRPC automatically generates responses based on your Protobuf definitions. Here you would normally have some application logic that you want to test so this code might live elsewhere. The randomly generated data might work in a few scenarios but in order to test

### 4. Defining Stub Responses
For more control over the responses you can define stubs. This allows you to simulate specific scenarios and test how your client handles different responses.

```go
container.MustAddStub(ctx, "connectrpc.eliza.v1.ElizaService/Say", &elizav1.SayResponse{
    Sentence: "I am setting this text!",
})
```
