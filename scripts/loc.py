#!/usr/bin/env python

"""
Calculates the size (in lines of code) of Sumatra source code
(excluding dependencies i.e. code not written by us).
"""

import os
pj = os.path.join

DIRS = [".",
        "tests",
        ]

g_long_fmt = False

def is_blacklisted(name):
    return False


def count_file(name):
    return name.endswith(".go") and not is_blacklisted(name)


def loc_for_file(filePath):
    loc = 0
    with open(filePath, "r") as f:
        for line in f:
            loc += 1
    return loc


def get_locs_for_dir(srcDir, dir):
    d = pj(srcDir, dir)
    files = os.listdir(dir)
    locs_per_file = {}
    for f in files:
        if not count_file(f):
            continue
        locs_per_file[f] = loc_for_file(pj(d, f))
    return locs_per_file


def get_dir_loc(locs_per_file):
    return sum(locs_per_file.values())


def short_format(locs_per_dir):
    loc_total = 0
    for dir in sorted(locs_per_dir.keys()):
        locs_per_file = locs_per_dir[dir]
        print("%6d %s" % (get_dir_loc(locs_per_file), dir))

        for file in sorted(locs_per_file.keys()):
            loc = locs_per_file[file]
            loc_total += loc
    print("\nTotal: %d" % loc_total)

def long_format(locs_per_dir):
    loc_total = 0
    for dir in sorted(locs_per_dir.keys()):
        locs_per_file = locs_per_dir[dir]
        print("\n%s: %6d " % (dir, get_dir_loc(locs_per_file)))

        for file in sorted(locs_per_file.keys()):
            loc = locs_per_file[file]
            print(" %-25s %d" % (file, loc))
            loc_total += loc
    print("\nTotal: %d" % loc_total)

def main():
    # we assume the script is run from top-level directory as
    # python ./script/loc.py
    srcDir = "."
    locs_per_dir = {}
    for dir in DIRS:
        locs_per_dir[dir] = get_locs_for_dir(srcDir, dir)

    if g_long_fmt:
        long_format(locs_per_dir)
    else:
        short_format(locs_per_dir)

if __name__ == "__main__":
    main()
