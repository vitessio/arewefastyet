#  Copyright 2021 The Vitess Authors.
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

from typing import Optional
from github import Github, GithubException


def resolve_ref(initial_ref: str):
    is_pr = False
    ref = get_sha_from_ref(initial_ref)
    if ref is None:
        pr_nb = int(initial_ref)
        ref = get_sha_from_pr(pr_nb)
        if ref is not None:
            is_pr = True
    return ref, is_pr


def get_sha_from_ref(ref) -> Optional[str]:
    c = Github()
    try:
        sha = c.get_repo("vitessio/vitess").get_commit(ref).sha
    except GithubException as e:
        print(e.data.get("message"))
        return None
    return sha


def get_sha_from_pr(pr: int) -> Optional[str]:
    c = Github()
    try:
        sha = c.get_repo("vitessio/vitess").get_pull(pr).head.sha
    except GithubException as e:
        print(e.data.get("message"))
        return None
    return sha
