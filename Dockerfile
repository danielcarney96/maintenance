FROM ubuntu:latest

ARG GH_EMAIL
ARG GH_NAME
ARG GH_AUTH_TOKEN

# Dependencies
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
RUN apt-get update -y -q
RUN apt-get install dialog apt-utils git openssh-client software-properties-common -y
RUN add-apt-repository ppa:ondrej/php
RUN apt install php-cli unzip curl nano gh -y

# SSH
RUN mkdir /root/.ssh/
RUN ln -s /run/secrets/id_rsa /root/.ssh/id_rsa
RUN echo "Host *.trabe.io\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN ssh-keyscan bitbucket.org >> /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

RUN mkdir repositories

# Node
RUN DEBIAN_FRONTEND='noninteractive' apt install nodejs npm -y
RUN curl https://raw.githubusercontent.com/creationix/nvm/master/install.sh | bash 

# Composer
RUN curl -sS https://getcomposer.org/installer -o /tmp/composer-setup.php
RUN php /tmp/composer-setup.php --install-dir=/usr/local/bin --filename=composer

# Configure git
RUN git config --global user.email $GH_EMAIL
RUN git config --global user.name $GH_NAME
RUN echo $GH_AUTH_TOKEN | gh auth login --with-token
