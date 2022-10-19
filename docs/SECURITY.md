<!--
Copyright (c) 2020 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
   
    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Security Policy

The CSM services/repositories are inspected for security vulnerabilities via [gosec](https://github.com/securego/gosec).

Every issue detected by `gosec` is mapped to a [CWE (Common Weakness Enumeration)](http://cwe.mitre.org/data/index.html) which describes in more generic terms the vulnerability.  The exact mapping can be found at https://github.com/securego/gosec in the issue.go file. The list of rules checked by `gosec` can be found [here](https://github.com/securego/gosec#available-rules).

In addition to this, there are various security checks that get executed against a branch when a pull request is created/updated.  Please refer to [pull request](/docs/CONTRIBUTING.md#pull-requests) for more information.

## Reporting a Vulnerability

Please report a vulnerability by opening an Issue in this repository.
