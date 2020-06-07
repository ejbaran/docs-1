#!/usr/bin/python3

import os
import re
import argparse

# Script to update redirect links throughout the docs
# Prints out bad links during the process
# Usage: 
# $ python3 update_links.py --replace <current_link> <new_link>
# $ python3 update_links.py

rootdir = os.path.abspath("../docs")

link_regex = r"\[(.+?)\]\(([^\(\)]+?)\)"

def traverse_file_tree(args):
    """
    Traverse through the documentation file structure to
    check for bad links. 

    If user specifies a replacement link, make replacements.
    """

    if args.replace:
        files_updated = 0
        links_updated = 0
        old_link = os.path.abspath(args.replace[0])
        new_link = os.path.abspath(args.replace[1])
        # check that new link is valid, old link likely is not
        if os.path.exists(new_link.split("#")[0]):
            for dirpath, dirnames, filenames in os.walk(rootdir):
                for filename in filenames:
                    if filename.endswith('.md'):
                        os.chdir(dirpath)
                        links_replaced = check_and_replace_links(os.path.join(dirpath, filename), (old_link, new_link))
                        if links_replaced > 0:
                            links_updated += links_replaced
                            files_updated += 1
            print("{} links updated in {} files.".format(links_updated, files_updated))
        else:
            print("New link ({}) does not exist. Please enter a valid replacement link.".format(new_link))

def check_and_replace_links(filename, replace=None):
    """
    Opens file, finds links, and makes replacements
    as specified by user. Will print out bad links it finds
    along the way.
    """

    openfile = open(filename, "r")
    text = openfile.read()
    openfile.close()

    matches = re.finditer(link_regex, text)
    adjust_index = 0
    links_replaced = 0
    for match in matches:
        name, old_link, hashtag = parse_link(match)

        # quick and dirty weed out of external links or links to headers within the same document
        if len(old_link) > 0 and not old_link.startswith("http"): 

            # catch any links that are not relevant but also not caught in external check above
            try:
                old_link_abs = os.path.abspath(old_link)
            except:
                old_link_abs = ""

            if replace:
                if old_link_abs == replace[0]:

                    # get the relative path based on the current directory
                    new_link = os.path.relpath(replace[1])

                    # if user specifies a new hashtag, use that one, otherwise use old hashtag
                    if "#" in new_link:
                        hashtag = "#" + new_link.split("#")[1]
                        new_link = new_link.split("#")[0]
                    final_str = "[{}]({}{})".format(name, new_link, hashtag)
                    print(final_str)
                    # update indices based on new link length
                    text = text[:match.start() + adjust_index] + final_str + text[match.end() + adjust_index:]
                    adjust_index += (len(final_str) - len(match.group()))
                    links_replaced += 1

            # goal cli files do not always reference the .md file. This causes an error in Github,
            # but is handled in production. TO DO: fix goal cli generation scripts so works both in testing
            # and on production servers. For now, just ignore to avoid clogging up output and not seeing *real*
            # bad links.
            if os.path.basename(old_link_abs).endswith(".md") and not os.path.exists(old_link_abs):
                print("Name: {}, Link: {}, Hashtag: {}".format(name, old_link, hashtag))
                print("Bad link: {} in filename: {}".format(old_link_abs, filename))
    if links_replaced > 0:
        openfile = open(filename, 'w')
        openfile.write(text)
        openfile.close()
    return links_replaced

def parse_link(match):
    """
    Return a 3-tuple of the text, the link, and the hashtag 
    given the match object.
    """
    name = match.group(1)
    link = match.group(2)
    if "#" in link:
        if link.startswith("#"):
            old_link = ""
            hashtag = "#" + link.split("#")[0]
        else:
            old_link = link.split("#")[0]
            hashtag = "#" + link.split("#")[1]
    else:
        old_link = link
        hashtag = ""
    return (name, old_link, hashtag)

parser = argparse.ArgumentParser(description="Checks for bad links in the docs repo and replace any links specified via --replace arg.")
parser.add_argument("--replace", type=str, nargs=2, 
                    default=None, required=False, help="Provide the old link, followed by the link to replace with.")
args = parser.parse_args()
traverse_file_tree(args)
