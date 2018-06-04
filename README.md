# Chronologist 🎞

Chronologist is a Kubernetes controller that syncs your Helm chart deployments 
with Grafana annotations.

Key features:

- For each Helm release you install/upgrade creates related Grafana annotation
- Determines if it was a rollout or rollback and tags annotation appropriately
- When you purge delete a release, deletes corresponding annotation

## Contributing

Contributions are welcome!

See [docs/development.md](docs/development.md) for detailed instructions on 
how to run development environment for Chronologist.
