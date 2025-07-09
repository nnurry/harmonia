import yaml
import json
import uuid
from dataclasses import dataclass, field


@dataclass
class UserData:
    hostname: str
    ssh_user: str
    ssh_public_keys_content: list[str]

    def to_ci_format(self) -> str:
        user_data_content = {
            "hostname": self.hostname,
            "manage_etc_hosts": True,
            "disable_root_pw": True,
            "users": [
                {
                    "name": self.ssh_user,
                    "sudo": "ALL=(ALL) NOPASSWD:ALL",
                    "ssh_authorized_keys": self.ssh_public_keys_content,
                }
            ],
        }

        return (
            "#cloud-config"
            + "\n"
            + yaml.dump(
                user_data_content, indent=2, default_flow_style=False, sort_keys=False
            )
        )


@dataclass
class NetworkConfig:
    ip_address: str
    nameservers: list[str]
    gateway: str
    mac_address: str

    def to_ci_format(self) -> str:
        network_config = {
            "network": {
                "version": 2,
                "ethernets": {
                    "eth0": {
                        "dhcp4": False,
                        "addresses": [f"{self.ip_address}/24"],
                        "gateway4": self.gateway,
                        "macaddress": self.mac_address,
                        "nameservers": {
                            "addresses": self.nameservers,
                        },
                    }
                },
            },
        }

        return yaml.dump(
            network_config, indent=2, default_flow_style=False, sort_keys=False
        )


@dataclass
class MetaData:
    hostname: str
    prefix: str = "monochromatic"
    instance_id: str = field(init=False)

    def __post_init__(self):
        self.instance_id = f"{self.prefix}{self.hostname}-{uuid.uuid4().hex[:8]}"

    def to_ci_format(self) -> str:
        meta_data = {
            "instance-id": self.instance_id,
            "local-hostname": self.hostname,
        }
        return json.dumps(meta_data, indent=2)


class CloudInitConfig:
    user_data: UserData
    meta_data: MetaData
    network_config: NetworkConfig

    # TODO: add required fields in the arguments
    def __init__(
        self,
        hostname: str,
        ssh_user: str,
        ssh_public_keys: list[str],
        ipv4_address: str,
        nameservers: list[str],
        ipv4_gateway_address: str,
        mac_address: str,
    ):
        self.user_data = UserData(
            hostname=hostname,
            ssh_user=ssh_user,
            ssh_public_keys_content=ssh_public_keys,
        )

        self.meta_data = MetaData(hostname=hostname)

        self.network_config = NetworkConfig(
            ip_address=ipv4_address,
            nameservers=nameservers,
            gateway=ipv4_gateway_address,
            mac_address=mac_address,
        )

    @staticmethod
    def from_dict(data: dict):
        return CloudInitConfig(
            hostname=data["hostname"],
            ssh_user=data["ssh_user"],
            ssh_public_keys=data["ssh_publickeys"],
            ipv4_address=data["ipv4_address"],
            nameservers=data["nameservers"],
            ipv4_gateway_address=data["ipv4_gateway_address"],
            mac_address=data["mac_address"],
        )
