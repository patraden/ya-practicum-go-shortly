# multichecker

## Overview

- Implements all standard SA* staticcheck [checks](https://staticcheck.dev/docs/checks/)
- Implements https://github.com/gostaticanalysis/emptycase linter
- Implements https://github.com/nishanths/exhaustive linter
- Checks that there is no direct call to `os.Exit` from main func of the main package


## How to use?
```bash
make staticlint:build
make staticlint:help
make staticlint
```