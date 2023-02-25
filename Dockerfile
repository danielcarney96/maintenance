FROM ubuntu:latest

# Dependencies
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections
RUN apt-get update -y -q
RUN apt-get install dialog apt-utils git openssh-client software-properties-common -y
RUN add-apt-repository ppa:ondrej/php
RUN apt install php-cli unzip curl nano -y

# SSH
RUN mkdir /root/.ssh/
COPY id_rsa /root/.ssh
RUN chmod 600 /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan bitbucket.org >> /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

# Node
RUN DEBIAN_FRONTEND='noninteractive' apt install nodejs npm -y
RUN curl https://raw.githubusercontent.com/creationix/nvm/master/install.sh | bash 

# Composer
RUN curl -sS https://getcomposer.org/installer -o /tmp/composer-setup.php
RUN php /tmp/composer-setup.php --install-dir=/usr/local/bin --filename=composer
