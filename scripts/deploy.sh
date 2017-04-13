#!/bin/bash

# Go to Parent Directory and Execute

echo `pwd`

if [[ $(git status -s) ]]
then
    echo "Commit or stash changes before deploy."
    exit 1;
fi

git checkout master
git pull

exit 0
