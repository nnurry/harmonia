import os
import subprocess
from config.cloud_init import CloudInitConfig
from datetime import datetime, timezone


def generate_ci_base_dir(ci_root_dir: str, prefix_name: str):
    ts_version = datetime.now(tz=timezone.utc).strftime("%Y%m%d_%H%M%S")
    return os.path.abspath((ci_root_dir + prefix_name + "/" + ts_version))


def write_cloud_init_iso_to_disk(
    config: CloudInitConfig,
    ci_base_dir: str,
    iso_filename: str = "cidata.json",
):
    os.makedirs(ci_base_dir, exist_ok=True)

    user_data_path = os.path.join(ci_base_dir, "user-data")
    meta_data_path = os.path.join(ci_base_dir, "meta-data")
    network_config_path = os.path.join(ci_base_dir, "network-config")

    with (
        open(user_data_path, "w") as user_data_file,
        open(meta_data_path, "w") as meta_data_file,
        open(network_config_path, "w") as network_config_file,
    ):
        user_data_file.write(config.user_data.to_ci_format())
        meta_data_file.write(config.meta_data.to_ci_format())
        network_config_file.write(config.network_config.to_ci_format())

    iso_path = os.path.join(ci_base_dir, iso_filename)
    cmd_parts = [
        "mkisofs",
        "-output",
        iso_path,
        "-volid",
        "cidata",
        "-joliet",
        "-r",
        user_data_path,
        meta_data_path,
        network_config_path,
    ]

    subprocess.run(cmd_parts, check=True, capture_output=True, text=True, shell=False)
    return iso_path
