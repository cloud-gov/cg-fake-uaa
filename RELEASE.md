Here's how to issue a new release.

1. Get [goxc][] and install it.
2. Get a GitHub API key by creating a new [personal access token][].
3. Make a `.goxc.local.json` file in the root directory of your
   repository checkout with the following content:

   ```json
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
5. Update `CHANGELOG.md` by moving the "Unreleased" section to a
   new section for your new version. Commit the changes with a
   message like "Bumped version to v1.0.4.". Push your changes to
   GitHub.
6. Tag your version and push it to GitHub. For instance, if you're
   releasing v1.0.4, do:

   ```
   git tag -a v1.0.4
   git push origin v1.0.4
   ```

7. Run `go generate`.
8. Run `goxc`. Note that you may need to symlink your project
   directory somewhere under your `GOPATH` and run it from there,
   which is weird.
9. You might want to visit your new [release][] and edit its
   description so it doesn't just say "built by goxc". Consider
   copy-pasting the section for your new version from `CHANGELOG.md`
   into the description.

[goxc]: https://github.com/laher/goxc
[personal access token]: https://github.com/settings/tokens
[release]: https://github.com/18F/cg-fake-uaa/releases
