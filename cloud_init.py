import os
import yaml
import json
import uuid
import subprocess
from datetime import datetime, timezone
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


@dataclass
class CloudInitConfig:
    user_data: UserData
    meta_data: MetaData
    network_config: NetworkConfig

    @staticmethod
    def from_dict(data: dict) -> "CloudInitConfig":
        user_data = UserData(
            data["hostname"],
            data["ssh_user"],
            data["ssh_publickeys"],
        )

        meta_data = MetaData(data["hostname"])

        network_config = NetworkConfig(
            data["ipv4_address"],
            data["nameservers"],
            data["ipv4_gateway_address"],
            data["mac_address"],
        )

        return CloudInitConfig(user_data, meta_data, network_config)

    def write_to_disk(self, base_dir: str):
        os.makedirs(base_dir, exist_ok=True)

        user_data_path = os.path.join(base_dir, "user-data")
        meta_data_path = os.path.join(base_dir, "meta-data")
        network_config_path = os.path.join(base_dir, "network-config")

        with (
            open(user_data_path, "w") as user_data_file,
            open(meta_data_path, "w") as meta_data_file,
            open(network_config_path, "w") as network_config_file,
        ):
            user_data_file.write(self.user_data.to_ci_format())
            meta_data_file.write(self.meta_data.to_ci_format())
            network_config_file.write(self.network_config.to_ci_format())


@dataclass
class CloudInitISO:
    cloud_init_config: CloudInitConfig
    vm_cloud_init_root_dir: str
    iso_filename: str = "cidata.json"
    vm_cloud_init_base_dir: str = field(init=False)
    ts_version: str = field(init=False)

    def __post_init__(self):
        self.ts_version = datetime.now(tz=timezone.utc).strftime("%Y%m%d_%H%M%S")
        self.vm_cloud_init_base_dir = os.path.abspath(
            (
                self.vm_cloud_init_root_dir
                + self.cloud_init_config.user_data.hostname
                + "/"
                + self.ts_version
            )
        )

    def write_to_disk(self):
        os.makedirs(self.vm_cloud_init_base_dir, exist_ok=True)

        self.cloud_init_config.write_to_disk(self.vm_cloud_init_base_dir)

        cmd_parts = [
            "mkisofs",
            "-output",
            os.path.join(self.vm_cloud_init_base_dir, self.iso_filename),
            "-volid",
            "cidata",
            "-joliet",
            "-r",
            os.path.join(self.vm_cloud_init_base_dir, "user-data"),
            os.path.join(self.vm_cloud_init_base_dir, "meta-data"),
            os.path.join(self.vm_cloud_init_base_dir, "network-config"),
        ]

        subprocess.run(
            cmd_parts, check=True, capture_output=True, text=True, shell=False
        )
