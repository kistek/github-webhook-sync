Clone git repos into folder based on description file/secret or crd

Repo description file in yaml form:

How to know which branch was used in a webhook? look in ref - refs/heads/branch

need security token for two factor auths
```
    git config credential.helper '!f() { sleep 1; echo "username=${GIT_USER}\npassword=${GIT_PASSWORD}"; }; f'
```

## Supported git Transfer Protocols
We currently support http because http is the only authentication method with scoped authorization

If requested I will add ssh and file support

## Future Work

* We will want to watch the secrets folder / config path and update accordingly so we don't need to do a rolling restart when the secret changes
* Allow target repo path to be configured per repository instead of en masse
* If there is sufficient use-case we could switch to CRD driven repository descriptions instead of secret ones
* What should we do if two repositories are named identically? Require different target path?
