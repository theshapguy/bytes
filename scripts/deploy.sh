#!/bin/bash

if [[ $(git status -s) ]]
then
    echo "Commit or stash changes before deploy."
    exit 1;
fi

git checkout master
git config user.email "no@email.com"
git config user.name "Shap Guy"
ssh-agent bash -c 'ssh-add ~/.ssh/bytes_id_rsa; git pull --rebase origin master'
supervisorctl restart nginx

exit 0
