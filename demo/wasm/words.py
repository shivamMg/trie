#!/usr/bin/env python

import string


WORDS_FILE = '/usr/share/dict/american-english'  # https://en.wikipedia.org/wiki/Words_(Unix)
TARGET_FILE = 'words.txt'


if __name__ == '__main__':
    with open(WORDS_FILE) as f:
        words = [line.rstrip('\n') for line in f.readlines()]
    words = [word for word in words if all([ch in string.ascii_lowercase for ch in word])]  # removes words: e.g. Abby (noun), accuser's (apostrophe)
    word_set = set(words)
    is_singular = lambda word: not word.endswith('s') or word[:-1] not in word_set
    words = [word for word in words if is_singular(word)]

    with open(TARGET_FILE, 'w') as f:
        for word in words:
            f.write(word + '\n')
    print(len(words), 'words written to', TARGET_FILE)