Clone git repos into folder based on description file/secret or crd

Repo description file in yaml form:

How to know which branch was used in a webhook? look in ref - refs/heads/branch

need security token for two factor auths
```
    git config credential.helper '!f() { sleep 1; echo "username=${GIT_USER}\npassword=${GIT_PASSWORD}"; }; f'
```

## Supported git Transfer Protocols
We currently support ssh, http, and file

## Future Work

We will want to watch the secrets folder and update accordingly so we don't need to do a rolling restart when the secret changes

If there is sufficient use-case we could switch to CRD driven repository descriptions instead of secret ones

If we get a webhook from a repo we don't know about, should we pull it? Thinking "no".
