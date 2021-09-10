<!--
Copyright (c) 2020 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
-->

# Security Policy

The CSM services/repositories are inspected for security vulnerabilities via [gosec](https://github.com/securego/gosec).

Every issue detected by `gosec` is mapped to a [CWE (Common Weakness Enumeration)](http://cwe.mitre.org/data/index.html) which describes in more generic terms the vulnerability.  The exact mapping can be found at https://github.com/securego/gosec in the issue.go file. The list of rules checked by `gosec` can be found [here](https://github.com/securego/gosec#available-rules).

In addition to this, there are various security checks that get executed against a branch when a pull request is created/updated.  Please refer to [pull request](/docs/CONTRIBUTING.md#pull-requests) for more information.

## Reporting a Vulnerability

Have you discovered a security vulnerability in this project?
We ask you to alert the maintainers by sending an email, describing the issue, impact, and fix - if applicable.

You can reach the CSM maintainers at karavi@dell.com.
