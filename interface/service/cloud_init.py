from abc import abstractmethod

from config.cloud_init import CloudInitConfig
from interface.service.base import IService


class ICloudInitService(IService):
    @abstractmethod
    def generate_ci_base_dir(self, ci_root_dir: str, prefix_name: str) -> str:
        raise NotImplementedError

    @abstractmethod
    def write_cloud_init_iso_to_disk(
        self,
        config: CloudInitConfig,
        ci_base_dir: str,
        iso_filename: str = "cidata.iso",
    ) -> str:
        return NotImplementedError
