# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.provision "shell", path: "bootstrap.sh"
  config.vm.network "forwarded_port", guest: 8080, host: 8080
end
