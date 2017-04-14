#!/bin/bash

if [[ $(git status -s) ]]
then
    echo "Commit or stash changes before deploy."
    exit 1;
fi

if [ "$TRAVIS_PULL_REQUEST" != "false" ]; then
      hugo -d public
      echo "Travis should not deploy from pull requests"
      exit 0
fi

# Generate site in public directory
hugo -d public
go build -o scripts/hooks scripts/hooks.go

# moving generated files to /tmp
rm -rf /tmp/public
mv public /tmp
mv scripts /tmp/public/scripts

git config user.name "Travis CI"
git config user.email "ci@travis.com"
git config --global push.default simple

git branch -a

git checkout source
git checkout -t -b master origin/master

# cleaning master branch
rm -rf `pwd`/*

# moving live website to clean website
mv /tmp/public/* .
rm -r /tmp/public/

#Adding and Committing with Date
git add -A

export TZ=":America/Denver"
now=$(date +"%Y-%m-%d %H:%M")

msg=":closed_book: Travis: rebuilding site $now"
if [ $# -eq 1 ]
  then msg="$1"
fi
#Automated Message
git commit -m "$msg"

# Push Branch Live
git push

# Go back to working branch
git checkout source
