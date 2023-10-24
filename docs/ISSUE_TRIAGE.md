<!--
Copyright (c) 2021-2022 Dell Inc., or its subsidiaries. All Rights Reserved.

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
# Triage issues

The main goal of issue triage is to categorize all incoming Dell Container Storage Modules issues and make sure each issue has all basic information needed for anyone else to understand and be able to start working on it.

> **Note:** This information is for Dell Container Storage Modules project Maintainers, Owners, and Admins. If you are a Contributor, then you will not be able to perform most of the tasks in this topic.

The core maintainers of the Dell Container Storage Modules project are responsible for categorizing all incoming issues and delegating any critical or important issue to other maintainers. Triage provides an important way to contribute to an open source project. 

Triage helps ensure issues resolve quickly by:

- Ensuring the issue's intent and purpose is conveyed precisely. This is necessary because it can be difficult for an issue to explain how an end user experiences a problem and what actions they took.
- Giving a contributor the information they need before they commit to resolving an issue.
- Lowering the issue count by preventing duplicate issues.
- Streamlining the development process by preventing duplicate discussions.

If you don't have the knowledge or time to code, consider helping with triage. The community will thank you for saving them time by spending some of yours.

## Labels

GitHub labels help organize and categorize GitHub issues.  The main categories of labels being used are:

- type/*: describes the type of the issue
- triage/*: describes the result of an issue triage
- area/*: the repository the issue is associated with
- Additional labels are applied to indicate when an issue needs triage and asking the community for help.

| Label | Description |
| - | - |
| type/bug | Something isn't working. This is the default label associated with a bug issue. |
| type/feature	| A feature. This label is applied to a feature issue. There is no template for a feature issue.  Maintainers can manually create a feature issue and assign it this label. Alternatively, if a feature request is triaged and determined to be important, this label can replace the type/feature-request label. |
| type/feature-request	| New feature request. This is the default label associated with a feature request issue. |
| type/question |	Ask a question. This is the default label associated with a question issue. |
| needs-triage	| Applied by default to bug and feature request issues so that the community understands triage is required. |
| triage/works-as-intended | Applied to a bug issue as part of triage when the issue works as intended by design. |
| triage/needs-information | Applied to an issue as part of triage indicating more information is required in order to work on it |
| triage/duplicate | Indicates the issue or pull request already exists. |
| triage/needs-investigation | Applied to an issue as part of triage indicating more investigation is required to reproduce and understand. Indicates that the issue cannot be easily reproduced and notifies contributors/maintainers that more investigation is required by anyone willing to contribute. |
| help wanted | Request for help from the community. |
| beginner friendly | The issue is suitable for a beginner to work on. |
| area/* | Used to associate the issue with a specific repository. |

## Find issues that need triage

The easiest way to find issues that haven't been triaged is to search for issues with the `needs-triage` label.

## Ensure the issue contains basic information

Make sure that the issue's author provided the standard issue information. The Dell Container Storage Modules project utilizes [GitHub issue templates](https://help.github.com/en/articles/creating-issue-templates-for-your-repository) to guide contributors to provide standard information that must be included for each type of template or type of issue.

### Standard issue information that must be included

The following section describes the various issue templates and the expected content.

#### Bug reports

Should explain what happened, what was expected and how to reproduce it together with any additional information that may help to give a complete picture of what happened such as screenshots, output and any environment related information that's applicable and/or maybe related to the reported problem:
- Dell Container Storage Modules version
- Bug description
- Expected behavior
- Logs
- Screenshots
- Platform & OS Dell Container Storage Modules is installed on
- Additional environment information

#### Feature requests

Should explain what feature that the author wants to be added and why that is needed.

#### Ask a Question requests

In general, if the issue description and title is perceived as a question no more information is needed.

### Good practices

To make it easier for everyone to understand and find issues they're searching for it's suggested as a general rule of thumbs to:

- Make sure that issue titles are named to explain the subject of the issue, has a correct spelling and doesn't include irrelevant information and/or sensitive information.
- Make sure that issue descriptions doesn't include irrelevant information.
- Make sure that issues do not contain sensitive information.
- Make sure that issues have all relevant fields filled in.
- Do your best effort to change title and description or request suggested changes by adding a comment.

> **Note:** Above rules are applicable to both new and existing issues.

### Do you have all the information needed to categorize an issue?

Depending on the issue, you might not feel all this information is needed. Use your best judgement. If you cannot triage an issue using what its author provided, explain kindly to the author that they must provide the above information to clarify the problem. Label issue with `triage/needs-information`.

If the author provides the standard information, but you are still unable to triage the issue, request additional information. Do this kindly and politely because you are asking for more of the author's time.  Label issue with `triage/needs-information`.

If the author does not respond to the requested information within the timespan of a week, close the issue with a kind note stating that the author can request for the issue to be reopened when the necessary information is provided.

If you receive a notification with additional information provided, but you are no long on issue triage, and you feel you do not have time to handle it, you should delegate it to the current person on issue triage.

## Categorizing an issue

### Duplicate issues

Make sure it's not a duplicate by searching existing issues using related terms from the issue title and description. If you think you know there is an existing issue, but can't find it, please reach out to one of the maintainers and ask for help. If you identify that the issue is a duplicate of an existing issue:

1. Add a comment `duplicate of #<issue number>`
2. Add the `triage/duplicate` label

### Bug reports

If it's not perfectly clear that it's an actual bug, quickly try to reproduce it.

**It's a bug/it can be reproduced:**

1. Add a comment describing detailed steps for how to reproduce it, if applicable.
2. If you know that maintainers won't be able to put any resources into it for some time then label the issue with `help wanted` and optionally `beginner friendly` together with pointers on which code to update to fix the bug. This should signal to the community that we would appreciate any help we can get to resolve this.

**It can't be reproduced:**
1. Either ask for more information needed to investigate it more thoroughly.  Provide details in a comment.
2. Either delegate further investigations to someone else.  Provide details in a comment.

**It works as intended/by design:**
1. Kindly and politely add a comment explaining briefly why we think it works as intended and close the issue.
2. Label the issue `triage/works-as-intended`.
3. Remove the `needs-triage` label.

**It does not work as intended/by design:**
1. Update the issue with additional details if needed
2. Remove the `needs-triage` label.
4. Assign the appropriate milestone

### Feature requests

1. If the feature request does not align with the product vision, add a comment indicating so, remove the `needs-triage` label and close the issue
2. Otherwise, add the appropriate labels and comments to the issue, remove the `needs-triage` label, and assign it to the correct milestone.

## Requesting help from the community

Depending on the issue and/or priority, it's always a good idea to consider signalling to the community that help from community is appreciated and needed in case an issue is not prioritized to be worked on by maintainers. Use your best judgement. In general, requesting help from the community means that a contribution has a good chance of getting accepted and merged.

In many cases the issue author or community as a whole is more suitable to contribute changes since they are experts in their domain. It's also quite common that someone has tried to get something to work using the documentation without success and made an effort to get it to work and/or reached out to the community to get the missing information.

1. Kindly and politely add a comment to signal to users subscribed to issue updates.
   - Explain that the issue would be nice to get resolved, but it isn't prioritized to work on by maintainers for an unforeseen future.
   - If possible or applicable, try to help contributors getting starting by adding pointers and references to what code/files need to be changed and/or ideas of a good way to solve/implement the issue.
2. Label the issue with `help wanted`.
3. If applicable, label the issue with `beginner friendly` to denote that the issue is suitable for a beginner to work on.

## Investigation of issues

When an issue has all basic information provided, but the reported problem cannot be reproduced at a first glance, the issue is labeled `triage/needs-information`. Depending on the perceived severity and/or number of [upvotes](https://help.github.com/en/articles/about-conversations-on-github#reacting-to-ideas-in-comments), the investigation will either be delegated to another maintainer for further investigation or put on hold until someone else (maintainer or contributor) picks it up and eventually starts investigating it.

Even if you don't have the time or knowledge to investigate an issue we highly recommend that you [upvote](https://help.github.com/en/articles/about-conversations-on-github#reacting-to-ideas-in-comments) the issue if you happen to have the same problem. If you have further details that may help to investigate the issue please provide as much information as possible.

## External PRs

Part of issue triage should also be triaging of external PRs. Main goal should be to make sure PRs from external contributors have an owner/reviewer and are not forgotten.

1. Check new external PRs which do not have a reviewer.
1. Maintainers need to ensure the pull request aligns with a GitHub bug or feature
1. If not and you know which issue it is solving, add the link yourself, otherwise ask the author to link the issue or create one.
1. Maintainers need to ensure the contribution is relevant and aligns with the product roadmap and priorities
1. Assign a reviewer based on who was handling the linked issue or what code or feature does the PR touches (look at who was the last to make changes there if all else fails).
1. Work with the contributor to guide them and help ensure our quality standards are met and that all GitHub checks pass. 

## GitHub Issue Management Workflow

The following section describes the triage workflow for new issues that are created.

### Bugs

This workflow starts off with a GitHub issue of type bug being created.

1. Collaborator or maintainer creates a GitHub bug using the appropriate GitHub issue template
2. By default the bug is assigned the type/bug and needs-triage labels

```                                                                                                                                                                                                                                                                                                              
                                               +--------------------------+                                                                              
                                               | New bug issue opened/more|                                                                              
                                               | information added        |                                                                              
                                               +-------------|------------+                                                                              
                                                             |                                                                                           
                                                             |                                                                                           
   +----------------------------------+  NO   +--------------|-------------+                                                                             
   | label: triage/needs-information  ---------  All required information  |                                                                             
   |                                  |       |  contained in issue?       |                                                                             
   +-----------------------------|----+       +--------------|-------------+                                                                             
                                 |                           | YES                                                                                       
                                 |                           |                                                                                           
   +--------------------------+  |                +---------------------+ YES +---------------------------------------+                                  
   |label:                    |  |                |  Duplicate Issue?   ------- Comment `Duplicate of #<issue number>`                                   
   |triage/needs-investigation|  | NO             |                     |     | Remove needs-triage label             |                                  
   +------|-------------------+  |                +----------|----------+     | label: triage/duplicate               |                                  
          |                      |                           | NO             +-----------------|---------------------+                                  
      YES |                      |                           |                                  |                                                        
          |      +---------------|----+   NO    +------------|------------+                     |                                                        
          |      |Needs investigation?|----------  Can it be reproduced?  |                     |                                                        
          |-------                    |         +------------|------------+                     |                                                        
                 +--------------------+                      | YES                              |                                                        
                                                             |                       +----------|----------+                                             
                                                +------------|------------+          |  Close Issue        |                                             
                 |-------------------------------  Works as intended?     |          |                     |                                             
                 |                     NO       |                         |          +----------|----------+                                             
                 |                              +------------|------------+                     |                                                        
                 |                                           |                                  |                                                        
                 |                                           | YES                              |                                                        
                 |                          +----------------|----------------+                 |                                                        
   +-------------|------------+             | Add comment                     |                 |                                                        
   |   Add area label         |             | Remove needs-triage label       ------------------|                                                        
   |   label: area/*          |             | label: triage/works-as-intended |                                                                          
   +-------------|------------+             +---------------------------------+                                                                          
                 |                                                                                                                                       
                 |                        +----------+                                                                                                   
                 |                        |   Done   ----------------------------------------                                                            
                 |                        +----|-----+                                      |                                                            
                 |                             |NO                                          |                                                            
                 |                             |                         +------------------|------------------+                                         
    +------------|-------------+          +----|----------------+ YES    |  Add details to issue               |                                         
    |Remove needs-triage label ------------  Signal Community?  ----------  label: help wanted                 |                                         
    |Assign milestone          |          |                     |        |  label: beginner friendly (optional)|                                         
    +--------------------------+          +---------------------+        +-------------------------------------+                                                                                                                                                                                                                                                                                                  
```
If the author does not respond to a request for more information within the timespan of a week, close the issue with a kind note stating that the author can request for the issue to be reopened when the necessary information is provided.

### Feature Requests

This workflow starts off with a GitHub issue of type feature request being created.

1. Collaborator or maintainer creates a GitHub `type/feature-request` using the appropriate GitHub issue template
2. By default the bug is assigned the `type/feature-request` and `needs-triage` labels

The following flow chart outlines the triage process

```         
                                            +---------------------------------+                                                  
                                            |New feature request issue opened/|                                                  
                                            |more information added           |                                                  
                                            +----------------|----------------+                                                  
                                                             |                                                                   
                                                             |                                                                   
    +---------------------------------+ NO     +-------------|------------+                                                      
    | label: triage/needs-information ---------- All required information |                                                      
    |                                 |        | contained in issue?      |                                                      
    +---------------------------------+        +-------------|------------+                                                      
                                                             |                                                                   
                                                             |                                                                   
    +---------------------------------------+                |                                                                   
    |Comment `Duplicate of #<issue number>` | YES +----------|----------+                                                        
    |Remove needs-triage label              -------  Duplicate issue?   |                                                        
    |label: triage/duplicate                |     |                     |                                                        
    +-----|---------------------------------+     +-----------|---------+                                                        
          |                                                   |NO                                                                
          |  +-------------------------+  NO   +-----------------------------+                                                   
          |  |Add comment              |--------  Is this a valid feature?   |                                                   
          |  |Remove needs-triage label|       |                             |                                                   
          |  +------|------------------+       +--------------|--------------+                                                   
          |         |                                         | YES                                                              
          |         |                                         |                                                                  
          |         |                         +---------------|--------------+                                                                                                    
          |         |                         | label: type/feature          |                                                   
        +-|---------|---+                     | Remove needs-triage label    |                                                   
        |  Close issue  |                     | Remove type/feature-request  |                                                   
        |               |                     | milestone?                   |                                                   
        +---------------+                     +------------------------------+
                                                              |
+-------------------------------------+      YES    +---------------------+
|  Add details to issue               ---------------  Signal Community?  |
|  label: help wanted                 |             |                     |
|  label: beginner friendly (optional)|             +---------------------+                            
+---|---------------------------------+                       | NO 
    |                                                         |
+------+                                                      |
| Done | ------------------------------------------------------                                        
+------+                 
```
If the author does not respond to a request for more information within the timespan of a week, close the issue with a kind note stating that the author can request for the issue to be reopened when the necessary information is provided.

In some cases you may receive a request you do not wish to accept.  Perhaps the request doesn't align with the project scope or vision.  It is important to tactfully handle contributions that don't meet the project standards.

1. Acknowledge the person behind the contribution and thank them for their interest and contribution
2. Explain why it didn't fit into the scope of the project or vision
3. Don't leave an unwanted contributions open.  Immediately close the contribution you do not wish to accept
