# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "generic/debian11"

  config.nfs.functional = false
  config.vm.synced_folder "./", "/vagrant", type: "virtiofs"

  config.vm.network "forwarded_port", guest: 80, host: 8084
  config.vm.provider "libvirt" do |vb|
    vb.memory = "2048"
    vb.cpus = "8"
    vb.memorybacking :access, :mode => "shared"
  end

  config.vm.provision "ansible" do |ansible|
    ansible.verbose = ""
    ansible.playbook = "provision/playbook.yml"
  end
end
