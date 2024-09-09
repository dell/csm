// Copyright © 2019-2022 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//go:generate go run semver/semver.go -f semver.tpl -o core_generated.go

package core

import "time"

var (
	// SemVer is the semantic version.
	SemVer = "unknown"

	// CommitSha7 is the short version of the commit hash from which
	// this program was built.
	CommitSha7 string

	// CommitSha32 is the long version of the commit hash from which
	// this program was built.
	CommitSha32 string

	// CommitTime is the commit timestamp of the commit from which
	// this program was built.
	CommitTime time.Time
)
