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
        "bench_cli.*"
    ]),
    python_requires='>=3.7',
)