#!/usr/bin/env python
__metaclass__ = type
from jujupy import (
    check_wordpress,
    Environment,
    format_listing,
    until_timeout,
)

from collections import defaultdict
import sys


def agent_update(environment, version):
    env = Environment(environment)
    for ignored in until_timeout(300):
        versions = env.get_status().get_agent_versions()
        if versions.keys() == [version]:
            break
        print format_listing(versions, version, environment)
        sys.stdout.flush()
    else:
        raise Exception('Some versions did not update.')


def main():
    try:
       agent_update(sys.argv[1], sys.argv[2])
    except Exception as e:
        print e
        sys.exit(1)


if __name__ == '__main__':
    main()
