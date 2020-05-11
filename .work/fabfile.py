import os 
from fabric.api import *
from fabric.task_utils import merge

env.roledefs = {
    'master': ['pi@donald'],
    'worker': ['pi@tick', 'pi@trick', 'pi@track']
}

#env.hosts = [ 'pi@tick', 'pi@trick', 'pi@track', 'pi@donald' ]
#env.hosts = [ 'cluster-tick', 'cluster-trick', 'cluster-track', 'cluster-donald' ]
#env.password = 'raspberry'
env.password = 'cluster'


@serial
@roles('master', 'worker')
def cluster_status():
    run('uname -n')
    run('neofetch')
    run('df -h')

@serial
@roles('master', 'worker')
def cluster_pieeprom():
    run('vcgencmd version')
    run('vcgencmd bootloader_config')
#    run('ls /lib/firmware/raspberrypi/bootloader/stable')
#    run('cp /lib/firmware/raspberrypi/bootloader/stable/pieeprom-2020-04-16.bin /tmp/pieeprom.bin')
#    run('ls -la /tmp/pi*')
#    run('rpi-eeprom-config /tmp/pieeprom.bin > /tmp/bootconf.txt')
#    run('ls -la /tmp/boot*')
#    run('cat /tmp/bootconf.txt')
#    run('sed -e \'s/WAKE_ON_GPIO=1/WAKE_ON_GPIO=0/g\' -e \'s/POWER_OFF_ON_HALT=0/POWER_OFF_ON_HALT=1/g\' /tmp/bootconf.txt >/tmp/bootconf.chg')
#    run('cat /tmp/bootconf.chg')
#    run('rpi-eeprom-config --out /tmp/pieeprom-new.bin --config /tmp/bootconf.chg /tmp/pieeprom.bin')
#    run('ls -la /tmp/pi*')
#    sudo('rpi-eeprom-update -d -f /tmp/pieeprom-new.bin')
#    run('vcgencmd version')

@serial
@roles('master', 'worker')
def cluster_pieeprom_reboot():
    sudo('reboot now')

@serial
@roles('worker')
def cluster_nfs_worker():
    run('tree /mnt/nfs')
    run('du -h /mnt/nfs')

@serial
@roles('master')
def cluster_nfs_master():
    run('tree /media/usb/nfsshare')
    run('du -h /media/usb/nfsshare')

#def cluster_nfs_status():
#    cluster_nfs_worker()
#    cluster_nfs_master()

@serial
@roles('master', 'worker')
def change_passwd(user, passwd):
    sudo('echo %s:%s | chpasswd' % (user, passwd))

@parallel
@roles('master', 'worker')
def post_install():
    sudo('raspi-config --expand-rootfs')
    sudo('raspi-config nonint do_change_locale de_DE.UTF-8')
    sudo('raspi-config nonint do_configure_keyboard de')
    update()
    install_basics()
    sudo('reboot now')

@parallel
@roles('master', 'worker')
def update():
    sudo('apt-get update')
    sudo('apt-get upgrade -y')
    sudo('apt-get autoremove -y')

@serial
@roles('master', 'worker')
def install_basics():
    sudo('apt-get install -y neofetch htop glances lsof tcpdump iotop iftop tmux')

@parallel
@roles('master', 'worker')
def cluster_down():
    sudo('shutdown now')

@parallel
@roles('master', 'worker')
def cluster_reboot():
    sudo('reboot now')

@parallel
@roles('master', 'worker')
def cmd(command):
    sudo(command)

@roles('naster')
def install_server_nfs():
    # USB-Stick in master top usb
    # https://raspberrytips.com/format-mount-usb-drive/
    # https://pimylifeup.com/raspberry-pi-nfs/
    # https://www.elektronik-kompendium.de/sites/raspberry-pi/2007061.htm
    # https://www.htpcguides.com/configure-nfs-server-and-nfs-client-raspberry-pi/
    # lsblk
    # sudo fdisk -l
    # sudo fdisk /dev/sda
    # sudo mkfs.ext4 /dev/sda1
    # ls /dev/disk/by-uuid
    # sudo vi /etc/fstab -> UUID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx /media/usb ext4 defaults 0
    sudo('mkdir -p /media/usb')
    sudo('mount -t ext4 -o defaults /dev/sda1 /media/usb')
    sudo('touch /media/usb/usbstick')
    sudo('mkdir -p /media/usb/nfsshare')
    sudo('chmod 755 /media/usb/nfsshare')
    sudo("chown pi:pi -R /media/usb")
    sudo('apt-get install nfs-kernel-server -y')
    # id pi -> uid=1000 pid=1000
    # sudo vi /etc/exports -> /media/usb/nfsshare *(rw,all_squash,insecure,async,no_subtree_check,anonuid=1000,anongid=1000)
    # sudo exportfs -ra
    # sudo service nfs-kernel-server restart
    # sudo service nfs-kernel-server status
    # sudo rpcinfo -p
    # hostname -I

@serial
@roles('worker')
def install_client_nfs():
    #sudo('apt-get install nfs-common -y')
    #sudo('mkdir -p /mnt/nfs')
    #sudo('chown -R pi:pi /mnt/nfs')
    # sudo mount 192.168.1.11:/media/usb/nfsshare /mnt/nfs
    # sudo vi /etc/fstab -> 192.168.1.11:/media/usb/nfsshare   /mnt/nfs   nfs    rw  0  0
    # sudo('echo "192.168.1.11:/media/usb/nfsshare   /mnt/nfs   nfs    rw  0  0" >> /etc/fstab')
    run('mkdir -p /mnt/nfs/hk')
    run('mkdir -p /mnt/nfs/logs')
    run('mkdir -p /mnt/nfs/archive')


