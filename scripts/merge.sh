#!/bin/bash

if [[ $(git status -s) ]]
then
    echo "Commit or stash changes before deploy."
    exit 1;
fi

if [ "$TRAVIS_PULL_REQUEST" != "false" ]; then
      echo "travis should not deploy from pull requests"
      exit 0
fi

# Generate site in public directory
hugo -d live/
ls public
go build -o scripts/hooks scripts/hooks.go

# moving generated files to /tmp
rm -rf /tmp/live
mv live /tmp
mv scripts /tmp/live/scripts
ls /tmp/live/

git branch -a

git checkout source
git checkout -t -b master origin/master

# cleaning master branch
rm -rf `pwd`/*

# moving live website to clean website
mv /tmp/live/* .
rm -r /tmp/live/

ls

#Adding and Committing with Date
git add -A

export TZ=":America/Denver"
now=$(date +"%Y-%m-%d %H:%M")

msg=":bookmark: travis: rebuilding site $now"
if [ $# -eq 1 ]
  then msg="$1"
fi
#Automated Message
git commit -m "$msg"

# Push Branch Live
git push

# Go back to working branch
git checkout source
