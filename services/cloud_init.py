import os
import subprocess
from config.cloud_init import CloudInitConfig
from datetime import datetime, timezone


def write_cloud_init_iso_to_disk(
    config: CloudInitConfig,
    ci_root_dir: str,
    iso_filename: str = "cidata.json",
):
    ts_version = datetime.now(tz=timezone.utc).strftime("%Y%m%d_%H%M%S")
    base_dir = os.path.abspath(
        (ci_root_dir + config.user_data.hostname + "/" + ts_version)
    )

    os.makedirs(base_dir, exist_ok=True)

    user_data_path = os.path.join(base_dir, "user-data")
    meta_data_path = os.path.join(base_dir, "meta-data")
    network_config_path = os.path.join(base_dir, "network-config")

    with (
        open(user_data_path, "w") as user_data_file,
        open(meta_data_path, "w") as meta_data_file,
        open(network_config_path, "w") as network_config_file,
    ):
        user_data_file.write(config.user_data.to_ci_format())
        meta_data_file.write(config.meta_data.to_ci_format())
        network_config_file.write(config.network_config.to_ci_format())

    cmd_parts = [
        "mkisofs",
        "-output",
        os.path.join(base_dir, iso_filename),
        "-volid",
        "cidata",
        "-joliet",
        "-r",
        user_data_path,
        meta_data_path,
        network_config_path,
    ]

    subprocess.run(cmd_parts, check=True, capture_output=True, text=True, shell=False)
