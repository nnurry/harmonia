from dataclasses import dataclass


@dataclass
class VirtualMachineConfig:
    name: str
    vcpu: int
    memory_gb: int
    disk_gb: int
