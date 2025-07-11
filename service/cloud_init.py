import os
from datetime import datetime, timezone

from interface.service.cli import IDiskImageService
from interface.service.cloud_init import ICloudInitService


class CloudInitService(ICloudInitService):
    def __init__(self, disk_img_svc: IDiskImageService):
        self.disk_img_svc = disk_img_svc

    def name(cls):
        return "cloud-init"

    def generate_ci_base_dir(self, ci_root_dir, prefix_name):
        ts_version = datetime.now(tz=timezone.utc).strftime("%Y%m%d_%H%M%S")
        return os.path.abspath(os.path.join(ci_root_dir, prefix_name, ts_version))

    def write_cloud_init_iso_to_disk(
        self, config, ci_base_dir, iso_filename="cidata.iso"
    ):
        os.makedirs(ci_base_dir, exist_ok=True)

        user_data_path = os.path.join(ci_base_dir, "user-data")
        meta_data_path = os.path.join(ci_base_dir, "meta-data")
        network_config_path = os.path.join(ci_base_dir, "network-config")
        iso_path = os.path.join(ci_base_dir, iso_filename)

        with (
            open(user_data_path, "w") as user_data_file,
            open(meta_data_path, "w") as meta_data_file,
            open(network_config_path, "w") as network_config_file,
        ):
            user_data_file.write(config.user_data.to_ci_format())
            meta_data_file.write(config.meta_data.to_ci_format())
            network_config_file.write(config.network_config.to_ci_format())

        self.disk_img_svc.create_cloud_init_iso(
            iso_path,
            [
                user_data_path,
                meta_data_path,
                network_config_path,
            ],
        )

        iso_path = os.path.join(ci_base_dir, iso_filename)
        return iso_path
