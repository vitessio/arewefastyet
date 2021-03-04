# Copyright 2021 The Vitess Authors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#    http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import setuptools

setuptools.setup(
    name="Vitess Benchmark",
    version="0.1",
    author="vitessio",
    url="https://github.com/vitessio/arewefastyet",
    classifiers=[
        "Programming Language :: Python :: 3",
        "Operating System :: OS Independent",
    ],
    entry_points={
        'console_scripts': [
            'clibench = bench_cli.cli:main',
        ],
    },
    packages=setuptools.find_packages(include=[
        "bench_cli",
        "bench_cli.*",
        "test",
        "test.*"
    ]),
    python_requires='>=3.7',
)