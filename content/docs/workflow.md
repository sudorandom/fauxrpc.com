---
title: 'Workflow'
weight: 15
slug: workflow
---

## Typical Workflow
Here's what a typical protobuf workflow might look like:

1. Define the schema
2. Generate Clients and Server Stubs
3. Implement your service
4. Deploy to a development environment
5. Tell the frontend engineers to start integrating with this new service.


## With FauxRPC
With FauxRPC, you can parallelize the last 3 steps. Here's what it might look like now:

1. Define the schema
2. Generate Clients and Server Stubs
3. The backend team implements the service *while the frontend engineers can already start working*

![](</workflow.svg>)

With FauxRPC, you can start iterating on the schema much quicker and in a much more natural way. Many times, flaws in the schema aren't discovered until you start using the service, so getting to a mocked-out prototype quicker can make the entire process faster.

In summation, you get:

- Parallel Development
- Faster Feedback Loops
- Reduced Integration Issues
- Simplified Debugging
- Reduced External Dependencies
