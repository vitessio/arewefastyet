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

import boto3
from typing import Optional
from botocore.exceptions import ClientError


def upload_file(file_name, bucket="arewefastyet", object_name=None, link_exp=86400) -> Optional[str]:
    if object_name is None:
        object_name = file_name

    s3_client = boto3.client('s3')
    try:
        s3_client.upload_file(file_name, bucket, object_name)
        url = s3_client.generate_presigned_url('get_object',
                                                    Params={'Bucket': bucket, 'Key': object_name},
                                                    ExpiresIn=link_exp
                                                    )
    except ClientError as e:
        print(e)
        return None
    return url
