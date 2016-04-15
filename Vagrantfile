$script = <<SCRIPT
sudo apt-get update
SCRIPT

Vagrant.configure("2") do |cluster|

    cluster.vm.define "rss" do |machine|
        machine.vm.box = "ubuntu/trusty64"
        machine.vm.hostname = "rss"

        machine.vm.provision "shell", inline: $script
        machine.vm.provision "docker", images: ["alpine:3.3", "golang:1.6"]

        machine.vm.provider "virtualbox" do |vbox|
            vbox.name = "rss"
            vbox.memory = 512
        end
    end

end
