# Support: Pre-Build Custom Procedures

Add your custom pre-build procedures to this directory. They will be executed in alphabetical order.

## Example

- Building a frontend before compiling the main binary.

```sh
# Tree structure
./project_root
└── support
    └── pre.d/                    # Custom pre-build procedures
        ├── 10-build-frontend.sh  # Example: Build frontend
        └── ...                   # More pre-build procedures... Just follow the same patterns
```

## Example procedures

1. [10-build-frontend](10-build-frontend.md)
2. 20-generate-assets.sh

   ```sh
   #!/usr/bin/env bash

   ## TITLE

   A short description of what this script does.

   ### DESCRIPTION

   A more detailed description of what this script does.

   ### AUTHOR

   The author of this script.

   ### VERSION

   The version of this script.

   ### DATE

   The date of this script.

   ## USAGE

   A short description of how to use this script.

   ## EXAMPLES

   A few examples of how to use this script.

   ## RULES

   1. The script should be executable.
   2. The script should be named in the format: `<number>-<description>.sh`
   3. The script should be in the format: `<number>-<description>.sh`
   4. The script should be in the format: `<number>-<description>.sh`
   ```

---

## Notes

- All the scripts inside this directory will be executed in alphabetical order
- We use the `make build` command to build the project, if you need to bypass this hook, use `RUN_PRE_SCRIPTS=false make build #or the make command you need...`
- Ensure that the scripts are executable: `chmod +x your_script_path`
- All the scripts should be executable directly, we don't use sourcing for this hook.
- All scripts are executed in a subshell, so changes to environment variables will not persist in the parent shell.
- All scripts run with `bash` shell.
- DO NOT use `sudo` or `root` privileges, will trigger an error.
- Use the `set` command as you wish, just be careful with it.

---

> Sweet coding **!** \
> If you like it, [let me know :)](https://github.com/kubex-ecosystem/xtui/stargazers).
