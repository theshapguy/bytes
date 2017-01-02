+++

title = "Automating Blog Hosting With Ansible"
description = "Provisioning Blog Server with Ansible"
date = "2017-01-03T12:32:17-07:00"
draft = false

+++

> **Micro Project #1**. Please read my [Hello World](http://theshapguy.com/post/hello-world) post to find out more about this blog. This fortnight I'm working with Ansible to automate blog hosting along with server management.
>
>
> This blog post is a hands down tutorial on `provisioning` a server by hosting a blog using Git and nginx with Ansible.

Ansible is a configuration management system. Ansible by default uses the SSH protocol to manage machines. There are other configuration tools like Puppet and Chef, however these require you to install a small piece of software on the client server. Hence, my preference to Ansible.

#### Installing Ansible

We only need to install Ansible on the local computer. However make sure the server has `Python 2.7` installed.

If you have [brew](http://brew.sh) installed on a Mac OS. You can simply do

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

#### Ansible Basics

There are two ways you can use Ansible depending on your use case. Firstly, use can use Ansible like `bash` scripts. These will be lists of `tasks` that you use to provision the server.

Secondly, we can use Ansible to create `roles`. Ansible uses roles to seperate it's blocks of tasks. This allows you to copy blocks from one playbook to another without any hassle. A `playbook` in essence is how your tasks and roles are executed. Addtionally, it can be uploaded into [Ansible Galaxy](https://galaxy.ansible.com) for other people to use it.


An Ansible Role needs to have a specific directory structure. You can manually create this directory structure; an easier way is to use `ansible-galaxy init newrole`.

This command creates the follwing directory strcuture.
```
.
├── defaults
│   └── main.yml
├── files
├── handlers
│   └── main.yml
├── meta
│   └── main.yml
├── tasks
│   └── main.yml
├── templates
└── vars
    └── main.yml
```
Now lets explain each of these.

**files**: This directory contains files that will be copied to the host.

**handlers**: These are basically tasks which are can be notified to be run after tasks are completed.

**meta**: This file cotains the meta information and dependency roles.

**templates**: This directory contains templates that use variable subsitution and, then copy it in your host.

**tasks**: It contains all of the tasks that the playbook needs to run.

**vars**: These are variables that you can use in your roles.

#### Prerequsites

1. Make sure you have access to a remote or local Ubuntu Server. I use [DigitalOcean](https://m.do.co/c/95bdc8dc8e65) to cater for all my server needs since it is quick and easy to spin up machines.
    **Quick Bonus**: if you sign up throught this [link](https://m.do.co/c/95bdc8dc8e65) you get $10 in  [DigitalOcean](https://m.do.co/c/95bdc8dc8e65) credit.

2. Also make sure you have a static website or a webpage as a git repository.

> All Ansible Playbooks are written in YAML.

Now lets get the ball rolling,

Clone the static website from the git respository. But, make sure you change the repo url. Additionally, make sure you set approriate variables in the vars folder. The `{{ user }}` variable has no default value so make sure you add one to the vars file.

```
- name: Clone Blog From Github using HTTPS
  git:
    repo: https://github.com/theshapguy/bytes.git
    dest: /home/{{ user }}/blog
    version: master
    update: yes
```

We will use nginx as a webserver. You can also use [Apache](https://httpd.apache.org/) as well, however the following configuration is for nginx.

```
- name: Install nginx
  apt:
    package: nginx
    state: latest
```

Now, lets create a `template` configuration for nginx virtual host.

```
server {
   listen 80;
   listen [::]:80;

   server_name {{ domain }}  www.{{ domain }};

    root /home/{{ user }}/blog;
    index index.html;

   location / {
       try_files $uri $uri/ =404;
   }
}
```
All `variables` are replaced using variables defined in the vars directory so make sure these are defined.

Copy this template and configure the nginx symlink. Also make sure that you have removed the default symlink which can cause confliting issues between two nginx configuration.

```
- name: Copy Nginx Config for Blog to Host
  template:
    src: blog.conf.http.j2
    dest: /etc/nginx/sites-available/blog.nginx.conf

- name: Add Symlink to Host
  file:
    src: /etc/nginx/sites-available/blog.nginx.conf
    dest: /etc/nginx/sites-enabled/blog
    state: link
  notify: restart nginx

- name: Remove default symlink
  file:
    path: /etc/nginx/sites-enabled/default
    state: absent
```

Now lets create a handler to restart nginx. If there are any changes to the file the `notify` tag is executed and nginx is `reloaded`.

```
- name: reload nginx
  service: name=nginx state=reloaded

  # basically you are doing -> sudo service nginx [options]
```

Now, your website should be running on your `{{ domain }}` as listed on your vars directory.

### Final Words

This is a quick and fast way to serve a website using Ansible. Please refer to my Github [repository](https://github.com/theshapguy/52WeeksOfCode/tree/master/blog-hosting-with-ansible) for the full playbook. The Github  [repository](https://github.com/theshapguy/52WeeksOfCode/tree/master/blog-hosting-with-ansible) also includes setup for a secure (https) website. It installs and renews `letsencrypt` certificate which makes your blog secure. Refer to this playbook for more tips and tricks that I use to host my [blog](theshapguy.com).


