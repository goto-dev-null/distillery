# Supported Sources

- GitHub
- GitLab
- Forgejo / Codeberg (Codeberg works out of the box; any Forgejo instance can be configured)
- Homebrew (binaries only, if anything has a dependency, it will not work at this time)
- Hashicorp (special handling for their releases, pointing to GitHub repos will automatically pass through)
- Kubernetes (special handling for their releases, pointing to GitHub repos will automatically pass through)

## Authentication

Distillery supports authentication for GitHub, GitLab, and Forgejo/Codeberg. There are CLI options to pass
in a token, but the preferred method is to set the appropriate environment variable using a tool like
[direnv](https://direnv.net/).

| Source | Environment variable | CLI flag |
|---|---|---|
| GitHub | `DISTILLERY_GITHUB_TOKEN` | `--github-token` |
| GitLab | `DISTILLERY_GITLAB_TOKEN` | `--gitlab-token` |
| Forgejo / Codeberg | `DISTILLERY_FORGEJO_TOKEN` | `--forgejo-token` |

This allows you to bypass any API rate limits that might be in place for unauthenticated requests, but more
importantly it allows you to install from private repositories that you have access to!

## Forgejo / Codeberg

See the dedicated [Forgejo / Codeberg](sources/forgejo.md) page for full usage and configuration details.