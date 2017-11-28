# vault-gitlab-deploy

Plugin for management of Gitlab deploy keys for [HashiCorp Vault](https://vault.io). Currently it provides the ability to add deploy keys to a Gitlab server with a limited TTL.

## Building

The plugin can be built as any other Golang project by using `go build` to produce the plugin binary.

## Deployment

Make sure that you have set the `plugin_directory` key in your Vault configuration. Move the binary to that directory.

Calculate the checksum of the binary, and add it to the plugin registry:

```
SHASUM=$(shasum -a 256 "/tmp/vault-plugins/vault-gitlab-deploy" | cut -d " " -f1)
vault write sys/plugins/catalog/example-plugin \
  sha_256="$SHASUM" \
  command="vault-gitlab-deploy"
```

Mount the plugin to your Vault:

```
vault mount -path=gitlab -plugin-name=gitlab plugin
```

## Usage

To add new deploy keys to a project in Gitlab you must first configure the plugin

```
vault write gitlab/config address=some.gitlab.com token=abCdaerkalfk deploy_key_ttl=500
```

To add a deploy key to the project specify the project id, if the deploy key should be able to push and the key.

```
vault write gitlab/100/deploykeys key="ssh-rsa AAAA..." can_push=false
```