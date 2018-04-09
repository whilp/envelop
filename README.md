# `envelop`

Run a command in an environment populated from 1password.

```
$ eval $(op signin my)
$ envelop GITHUB_USERNAME=github.username -- /bin/sh -c 'echo $GITHUB_USERNAME'
whilp
```

https://support.1password.com/command-line/
