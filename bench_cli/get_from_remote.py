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

import os
import pysftp


def get_from_remote(ip, username, src_dir, dest_dir, is_directory=False, create_dest=False):
    cnopts = pysftp.CnOpts()
    cnopts.hostkeys = None
    if create_dest is True and os.path.exists(dest_dir) is False:
        os.mkdir(dest_dir)
    with pysftp.Connection(host=ip, username=username, cnopts=cnopts) as sftp:
        if is_directory is True:
            sftp.get_r(src_dir, dest_dir, preserve_mtime=True)
        else:
            sftp.get(src_dir, dest_dir, preserve_mtime=True)
