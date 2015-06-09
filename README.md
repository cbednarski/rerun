# rerun

`rerun` is a small utility that reruns a command and writes stdout/stderr to a file whenever the exit code is non-zero.

`rerun` is intended to help identify tests that fail *sometimes*, and other quirky behavior. It automatically tracks failures and collects some basic statistics.

## Installation

    go install github.com/cbednarski/rerun

## Usage

    rerun command
