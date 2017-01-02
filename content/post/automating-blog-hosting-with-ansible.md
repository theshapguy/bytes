+++
title = "Automating Blog Hosting With Ansible"
description = ""
date = "2017-01-01T22:32:17-07:00"

+++

Ansible is a configuration management system. Ansible by default uses the SSH protocol manage machines.There are other configuration tools like Puppet and Chef, however most of them require you to install a small piece of software on the client server.

This blog post is a hands down tutorial on Ansible where it shows how to provision a server by hosting a blog on a server using Git and nginx.

Additionally, Ansible has new vocabularies. Please refer to [Ansible Glossary Docs](http://docs.ansible.com/ansible/glossary.html) if you do not understand some of the terms used.

This blog post will show you how to setup a DigitalOcean Droplet and automate your static website blog.


### How to install Ansible?

Since we talked about this above, we only need to install Ansible on the central server.

If you have [brew](http://brew.sh) installed on a Mac OS. You can just do

```
$ brew install ansible
```

For Ubuntu,

```

$ sudo apt-get install software-properties-common
$ sudo apt-add-repository ppa:ansible/ansible
$ sudo apt-get update
$ sudo apt-get install ansible

```
Please refer to [Ansible Installation Docs](http://docs.ansible.com/ansible/intro_installation.html#latest-releases-via-apt-ubuntu) to install on more OS's

### Step 0.1

Ansible uses Roles to seperate it's blocks of installation script. This allows you to copy blocks from one playbook (installation script) to another without any hassle. Addtionally, it can be uploaded into Ansible Galaxy for other people to use it.

```
files: This directory contains regular files that need to be transferred to the hosts you are configuring for this role. This may also include script files to run.
handlers: All handlers that were in your playbook previously can now be added into this directory.
meta: This directory can contain files that establish role dependencies. You can list roles that must be applied before the current role can work correctly.
templates: You can place all files that use variables to substitute information during creation in this directory.
tasks: This directory contains all of the tasks that would normally be in a playbook. These can reference files and templates contained in their respective directories without using a path.
vars: Variables for the roles can be specified in this directory and used in your configuration files.
```

Ansible uses roles to seperate blocks of tasks. This allows roles to be moved from one playbook to another. A role usually has a structure of the tree diagram below.

```
├── roles
│   ├── common
│   │   ├── files
│   │   ├── handlers
│   │   ├── meta
│   │   ├── tasks
│   │   ├── templates
│   │   └── vars
│
```



However a easier way is to do ```ansible-galaxy init git``` which makes it easier and we can also upload this into ansible galaxy for other people to use it. As we can see ansible-galaxy also creates a test directory.






### About this Article


