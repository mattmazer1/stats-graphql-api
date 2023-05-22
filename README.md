## Graphql api

A Graphql api to perform crud operations on current season stats for a roster of basketball players. All crud operations for players can only be done by an authenticated user with a valid JWT token. Subscriptions are implemented with pub sub from NATS to make it easy for a frontend to subscribe to updates. This api is deployed to AWS ECS.

stack - go, graphql, nats, postgres, docker, aws ecs
