# Forgejo / Codeberg

[Forgejo](https://forgejo.org) is an open-source, self-hostable Git forge. [Codeberg](https://codeberg.org)
is the most widely used public instance of Forgejo. Distillery supports Codeberg out of the box and can be
configured to work with any other Forgejo instance.

## Codeberg

Codeberg is supported via the `codeberg` prefix. No configuration is required.

### Examples

Install the latest release of a tool from Codeberg:

```console
dist install codeberg/owner/repo
```

Install a specific version:

```console
dist install codeberg/owner/repo@1.2.3
```

## Authentication

For private repositories or to avoid API rate limits, set a Forgejo/Codeberg access token.

The preferred method is via environment variable:

```bash
export DISTILLERY_FORGEJO_TOKEN=your_token_here
```

Or pass it directly on the command line:

```console
dist install --forgejo-token=your_token_here codeberg/owner/repo
```

Tokens can be created in your Codeberg account under **Settings → Applications → Access Tokens**.
The token only needs **read** access to repositories and releases.

## Other Forgejo Instances

Any self-hosted Forgejo (or Gitea-compatible) instance can be added as a named provider in your
configuration file. The provider name you choose becomes the routing prefix used on the command line.

=== "YAML"

    ```yaml
    providers:
      myforgejo:
        provider: forgejo
        base_url: https://git.example.com/api/v1
    ```

=== "TOML"

    ```toml
    [providers.myforgejo]
    provider = "forgejo"
    base_url = "https://git.example.com/api/v1"
    ```

Once configured, use the provider name as the prefix:

```console
dist install myforgejo/owner/repo
dist install myforgejo/owner/repo@2.0.0
```

The same `--forgejo-token` / `DISTILLERY_FORGEJO_TOKEN` is used for authentication across all
configured Forgejo providers.

## Directory Structure

Binaries installed from Codeberg are stored under the `codeberg` namespace:

```
~/.distillery/opt/codeberg/owner/repo/version/binary
```

Binaries installed from a custom configured Forgejo provider are stored under the provider's name as
defined in your configuration file:

```
~/.distillery/opt/myforgejo/owner/repo/version/binary
```

This means `dist list` will show the provider name, preserving a clear record of where each binary
came from.

## Aliases

You can create aliases for frequently used Forgejo repositories in your configuration file:

=== "YAML"

    ```yaml
    aliases:
      my-tool: codeberg/owner/repo
      pinned-tool: codeberg/owner/other-repo@1.0.0
    ```

=== "TOML"

    ```toml
    [aliases]
    my-tool = "codeberg/owner/repo"
    pinned-tool = "codeberg/owner/other-repo@1.0.0"
    ```
