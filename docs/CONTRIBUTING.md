<!--
Copyright (c) 2020 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
-->

# How to Contribute

Become one of the contributors to this project! We thrive to build a welcoming and open community for anyone who wants to use the project or contribute to it. There are just a few small guidelines you need to follow. To help us create a safe and positive community experience for all, we require all participants to adhere to the [Code of Conduct](CODE_OF_CONDUCT.md).

## Table of Contents

* [Become a contributor](#Become-a-contributor)
* [Submitting issues](#Submitting-issues)
* [Triage issues](#Triage-issues)
* [Your first contribution](#Your-first-contribution)
* [Branching](#Branching)
* [Signing your commits](#Signing-your-commits)
* [Pull requests](#Pull-requests)
* [Code reviews](#Code-reviews)
* [TODOs in the code](#TODOs-in-the-code)

## Become a contributor

You can contribute to this project in several ways. Here are some examples:

* Contribute to the CSM documentation and codebase.
* Report and triage bugs.
* Feature requests
* Write technical documentation and blog posts, for users and contributors.
* Help others by answering questions about this project.

## Submitting issues

All issues related to CSM, regardless of the service/repository the issue belongs to (see table above), should be submitted [here](https://github.com/dell/csm/issues). Issues will be triaged and labels will be used to indicate the type of issue. This section outlines the types of issues that can be submitted.  

### Report bugs

We aim to track and document everything related to CSM via the Issues page. The code and documentation are released with no warranties or SLAs and are intended to be supported through a community driven process.

Before submitting a new issue, make sure someone hasn't already reported the problem. Look through the [existing issues](https://github.com/dell/csm/issues) for similar issues.

Report a bug by submitting a [bug report](https://github.com/dell/csm/issues/new?labels=type%2Fbug%2C+needs-triage&template=bug_report.md&title=%5BBUG%5D%3A). Make sure that you provide as much information as possible on how to reproduce the bug.

When opening a Bug please include the following information to help with debugging:

1. Version of relevant software: this software, Kubernetes, Dell Storage Platform, Helm, etc.
2. Details of the issue explaining the problem: what, when, where
3. The expected outcome that was not met (if any)
4. Supporting troubleshooting information. __Note: Do not provide private company information that could compromise your company's security.__

An Issue __must__ be created before submitting any pull request. Any pull request that is created should be linked to an Issue.

### Feature request

If you have an idea of how to improve this project, submit a [feature request](https://github.com/dell/csm/issues/new?labels=type%2Ffeature-request%2C+needs-triage&template=feature_request.md&title=%5BFEATURE%5D%3A).

### Answering questions

If you have a question and you can't find the answer in the documentation or issues, the next step is to submit a [question.](https://github.com/dell/csm/issues/new?labels=type%2Fquestion&template=ask-a-question.md&title=%5BQUESTION%5D%3A)

We'd love your help answering questions being asked by other CSM users.

## Triage issues

Triage helps ensure that issues resolve quickly by:

* Ensuring the issue's intent and purpose is conveyed precisely. This is necessary because it can be difficult for an issue to explain how an end user experiences a problem and what actions they took.
* Giving a contributor the information they need before they commit to resolving an issue.
* Lowering the issue count by preventing duplicate issues.
* Streamlining the development process by preventing duplicate discussions.

If you don't have the knowledge or time to code, consider helping with _issue triage_. The CSM community will thank you for saving them time by spending some of yours.

Read more about the ways you can [Triage issues](ISSUE_TRIAGE.md).

## Your first contribution

Unsure where to begin contributing? Start by browsing issues labeled `beginner friendly` or `help wanted`.

* [Beginner-friendly](https://github.com/dell/csm/issues?q=is%3Aopen+is%3Aissue+label%3A%22beginner+friendly%22) issues are generally straightforward to complete.
* [Help wanted](https://github.com/dell/csm/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22) issues are problems we would like the community to help us with regardless of complexity.

When you're ready to contribute, it's time to create a pull request.

## Branching

* [Branching Strategy for CSM](BRANCHING.md)

## Signing your commits

We require that developers sign off their commits to certify that they have permission to contribute the code in a pull request. This way of certifying is commonly known as the [Developer Certificate of Origin (DCO)](https://developercertificate.org/). We encourage all contributors to read the DCO text before signing a commit and making contributions.

GitHub will prevent a pull request from being merged if there are any unsigned commits.

### Signing a commit

GPG (GNU Privacy Guard) will be used to sign commits.  Follow the instructions [here](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/signing-commits) to create a GPG key and configure your GitHub account to use that key.

Make sure you have your user name and e-mail set.  This will be required for your signed commit to be properly verified.  Check the following references:

* Setting up your github user name [reference](https://help.github.com/articles/setting-your-username-in-git/)
* Setting up your e-mail address [reference](https://help.github.com/articles/setting-your-commit-email-address-in-git/)

Once Git and your GitHub account have been properly configured, you can add the -S flag to the git commits:

```console
$ git commit -S -m your commit message
# Creates a signed commit
```

### Commit message format

CSM uses the guidelines for commit messages outlined in [How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/)

## Pull Requests

If this is your first time contributing to an open-source project on GitHub, make sure you read about [Creating a pull request](https://help.github.com/en/articles/creating-a-pull-request).

A pull request must always link to at least one GitHub issue. If that is not the case, create a GitHub issue and link it.

To increase the chance of having your pull request accepted, make sure your pull request follows these guidelines:

* Title and description matches the implementation.
* Commits within the pull request follow the formatting guidelines.
* The pull request closes one related issue.
* The pull request contains necessary tests that verify the intended behavior.
* If your pull request has conflicts, rebase your branch onto the main branch.

If the pull request fixes a bug:

* The pull request description must include `Fixes #<issue number>`.
* To avoid regressions, the pull request should include tests that replicate the fixed bug.

The CSM team _squashes_ all commits into one when we accept a pull request. The title of the pull request becomes the subject line of the squashed commit message. We still encourage contributors to write informative commit messages, as they becomes a part of the Git commit body.

We use the pull request title when we generate change logs for releases. As such, we strive to make the title as informative as possible.

Make sure that the title for your pull request uses the same format as the subject line in the commit message.

### Quality Gates for pull requests

GitHub Actions are used to enforce quality gates when a pull request is created or when any commit is made to the pull request. These GitHub Actions enforce our minimum code quality requirement for any code that get checked into the CSM Go code repository. If any of the quality gates fail, it is expected that the contributor will look into the check log, understand the problem and resolve the issue. If help is needed, please feel free to reach out the maintainers of the project for [support](SUPPORT.md).

#### Security scans

* [Golang Security Checker](https://github.com/securego/gosec) inspects source code for security vulnerabilities by scanning the Go AST.
* [Malware Scanner](https://github.com/dell/common-github-actions/tree/main/malware-scanner) inspects source code for malware.
* [Container Scanner](https://github.com/Azure/container-scan) scans containers for security vulnerabilities.

#### Code vetting

[GitHub action](https://github.com/dell/common-github-actions/tree/main/go-code-formatter-linter-vetter) that analyzes source code to report suspicious constructs such as Printf calls whose arguments do not align with the format string, abnormal or not used code in pull requests. Please refer to [vet](https://golang.org/cmd/vet/) for more information.

#### Code linting

[GitHub action](https://github.com/dell/common-github-actions/tree/main/go-code-formatter-linter-vetter) that analyzes source code to flag programming errors, stylistics errors, and suspicious constructs. Please refer to [Go lint](https://github.com/golang/lint) for more information.

#### Code formatting

[GitHub action](https://github.com/dell/common-github-actions/tree/main/go-code-formatter-linter-vetter) that analyzes source code to flag formatting errors. Please refer to [gofmt](https://golang.org/cmd/gofmt/) for more information.

#### Code sanitization

[GitHub action](https://github.com/dell/common-github-actions/tree/main/code-sanitizer) that analyzes source code for non-inclusive words and language.

#### Code build/test/coverage

[GitHub action](https://github.com/dell/common-github-actions/tree/main/go-code-tester) that runs Go unit tests and checks that the code coverage of each package meets a configured threshold (currently 90%). An error is flagged if a given pull request does not meet the test coverage threshold and blocks the pull request from being merged.

## Code Reviews

All submissions, including submissions by project members, require review. We use GitHub pull requests for this purpose. Consult [GitHub Help](https://help.github.com/articles/about-pull-requests/) for more information on using pull requests.

A pull request must satisfy following for it to be merged:

* A pull request will require at least 2 maintainer approvals.
* Maintainers must perform a review to ensure the changes adhere to guidelines laid out in this document.
* If any commits are made after the PR has been approved, the PR approval will automatically be removed and the above process must happen again.

## Code Style

For the Go code in the CSM repository, we expect the code styling outlined in [Effective Go](https://golang.org/doc/effective_go.html). In addition to this, we have the following supplements:

### Handle Errors

See [Effective Go](https://golang.org/doc/effective_go.html#errors) for details on handling errors.

Do not discard errors using _ variables. If a function returns an error, check it to make sure the function succeeded.  Handle the error, return it, or, in truly exceptional situations, panic.  This can be checked using the errcheck tool if you have it installed locally.

Do not log the error if it will also be logged by a caller higher up the call stack;  doing so causes the logs to become repetitive.  Instead, consider wrapping the error in order to provide more detail.  To see practical examples of this, see this bad practice and this preferred practice:

#### Bad

```go
package main

import (
    "errors"
    "log"
)

func main() {
    err := foo()
    if err != nil {
        log.Printf("error: %+v", err)
        return
    }
}

func foo() error {
    err := bar()
    if err != nil {
        log.Printf("error: %+v", err)
        return err
    }
    return nil
}

func bar() error {
    return errors.New("something bad happened")
}
```

#### Preferred

```go
package main

import (
    "errors"
    "fmt"
    "log"
)

func main() {
    err := foo()
    if err != nil {
        log.Printf("error: %+v", err)
        return
    }
}

func foo() error {
    err := bar()
    if err != nil {
        return fmt.Errorf("calling bar: %w", err)
    }
    return nil
}

func bar() error {
    return errors.New("something bad happened")
}
```

Do not use the github.com/pkg/errors package as it is now in maintenance mode since Go 1.13+ added official support for error wrapping. See [go1.13-errors](https://blog.golang.org/go1.13-errors) and [errwrap](https://github.com/fatih/errwrap) for more information.

### Gofmt

Run gofmt on your code to automatically fix the majority of mechanical style issues. Almost all Go code in the wild uses gofmt. The rest of this document addresses non-mechanical style points.

An alternative is to use goimports, a superset of gofmt which additionally adds (and removes) import lines as necessary.

A recommended approach is to ensure your editor supports running of goimports automatically on save.

### TODOs in the code

We don't like TODOs in the code or documentation. It is really best if you sort out all issues you can see with the changes before we check the changes in.
