# Harmonia: A Simple Hypervisor Manager

Harmonia is a small-scale hypervisor manager built for personal use. It's designed to be used with self-managed bare-metal nodes, providing a simple way to create and manage fleets of virtual machines using self-defined YAML/JSON configuration files.

**Disclaimer**: This tool is not intended for production environments.

### Features
- Create and delete virtual machine fleets.
- Configure VMs using YAML or JSON files.
- Built solely on Libvirt and SSH.

### Example Configuration

This is a sample YAML configuration file for a small VM fleet.

It demonstrates how to define a shared configuration for SSH and hypervisor connections, as well as specific settings for each individual virtual machine, including name, IP address, and resource allocation.

```
shared_config:
  general:
    base_vm_name: "base-VM"
  ssh:
    user: root
    authorized_key_contents:
      - ssh-ed25519 ABC XYZ
  cloud_init:
    nameservers:
      - "8.8.8.8"
      - "8.8.4.4"
    disable_root_pw: true
  hypervisor_connection:
    is_local_shell: false
    libvirt:
      connection_url: "qemu+ssh://root@hypervisor/system"
      keyfile_path: "/root/.ssh/hypervisor-id_ed25519"
    ssh:
      user: root
      host: hypervisor
      port: 22
      hostkey_callback_name: InsecureIgnoreHostKey
      privkey_auth_config:
        path: "/root/.ssh/hypervisor-id_ed25519"

virtual_machines:
  - name: "master-1"
    ip_address: "192.168.10.101"
    gateway_address: "192.168.10.1"
    vcpu: 2
    memory_gb: 8
    disk_gb: 50
    mac_address: "52:54:00:00:00:01"
    is_cow_clone: true

  - name: "worker-1"
    ip_address: "192.168.10.111"
    gateway_address: "192.168.10.1"
    vcpu: 2
    memory_gb: 16
    disk_gb: 50
    mac_address: "52:54:00:00:00:11"
    is_cow_clone: true
```

## RELEASE
- Version 0.0.0.1:
    - This version establishes the core functionality of creating and deleting virtual machine fleets on bare-metal nodes using configuration files.
