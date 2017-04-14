#!/bin/bash

if [[ $(git status -s) ]]
then
    echo "Commit or stash changes before deploy."
    exit 1;
fi

git checkout master
ssh-agent bash -c 'ssh-add ~/.ssh/bytes_id_rsa; git pull'
supervisorctl restart nginx

exit 0
