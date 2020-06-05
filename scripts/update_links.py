#!/usr/bin/python3

import os
import re

# Script to update redirect links throughout the docs
# Prints out bad links during the process
# Usage: 
# $ python3 update_links <current_link> <new_link>


# 1. recurse through files
# 2. get current filename and path
# open file and read contents
#
rootdir = os.path.abspath("../docs")

link_regex = r"\[(.+?)\]\(([^\(\)]+?)\)"

def check_and_replace_links(filename):
    """
    Opens file and checks links
    """
    os.chdir(rootdir)
    if os.path.exists(filename):
        openfile = open(filename, "r")
        text = openfile.read()
        openfile.close()
        print(text)
        matches = re.finditer(link_regex, text)
        for match in matches:
            linkify(match)
        #print(os.path.relpath(link, filename))
    else:
        print(filename)

def linkify(match):
    print(match.span())

def reformat_link():
    pass

def parse_header():
    header = "\#[^\]\[]\#]$"

def distinguish_link_type():
    pass


if __name__ == "__main__":
    check_and_replace_links('./docs/features/asa.md')
