# Monzo CLI

CLI for interacting with the Monzo APIs. Utilizes the
[Go Monzo API Client](../../).

## Command Docs

The command documentation can be found in the [docs/](docs/monzo.md)
subdirectory.

## Basic Usage

Here are some samples on basic usage of the CLI.

### Login

You can login/authenticate against the Monzo API with either a static access
token, or by using the OAuth2 authorization flow with a registered OAuth2
client.

If you wish to fetch an accounts entire transaction history, you will need to
use the OAuth2 login flow, as transaction data is restricted to the last 90
days 5 minutes after the access token is granted access to the account.

To login with a single token, run:

```shell
monzo login --access-token eyJ...
```

To login with the OAuth2 flow, run:

```shell
monzo login --client-id oauth2client_... --client-secret mnzconf...
```

Instead of CLI flags, you can also use environment variables:

* `MONZO_ACCESS_TOKEN`
* `MONZO_CLIENT_ID`
* `MONZO_CLIENT_SECRET`

Once you are authenticated, your access token and OAuth2 client details are
stored on disk in a token cache file. See [Caches](#caches) below.

### Logout

Logging out will delete the cache directory from disk and also attempt to
revoke the access token.

```shell
monzo logout
```

## Caches

The Monzo CLI stores certain persistent data on disk for use between
invocations. By default, this is stored in `$HOME/.monzo/`.

The directory can be overriden using the `MONZO_HOME_DIR` environment variable.

Cache files are unencrypted, so care should be taken to appropriately protect
the directory. By default the directory is created with `0600` permissions.
