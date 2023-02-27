# kubectl-clean-get

A tool to have a clean yaml output for kubectl get command.
It removes fields that are instance specific.
- uid
- selfLink
- resourceVersion
- creationTimestamp
- namespace

It can be use handly in CI to copy resources accross namespaces
