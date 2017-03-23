Here's how to issue a new release.

1. Get [goxc][] and install it.
2. Get a GitHub API key by creating a new [personal access token][].
3. Make a `.goxc.local.json` file in the root directory of your
   repository checkout with the following content:

```
{
        "Tasks": [
                "default",
                "publish-github"
        ],
        "TaskSettings": {
                "publish-github": {
                        "apikey": "<PUT YOUR API KEY HERE>"
                }
        },
        "ConfigVersion": "0.9"
}
```

4. Edit `.goxc.json` and increment the `PackageVersion`.
5. Run `go generate`.
6. Run `goxc`.

[goxc]: https://github.com/laher/goxc
[personal access token]: https://github.com/settings/tokens
