# rerun

rerun is a small utility that reruns a command and writes the output to a file whenever the exit status is non-zero. It's intended to help identify tests that fail *sometimes*, and identify other quirky behaviors in distributed systems or applications with non-deterministic scheduling (e.g. threading). rerun will track failures automatically and continue running so you can catch more than one type of periodic failure, and also get a sense for frequency.

## Usage

    rerun command
