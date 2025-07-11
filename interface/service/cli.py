from abc import abstractmethod
from subprocess import CompletedProcess

from interface.service.base import IService


class ICLIService(IService):
    @abstractmethod
    def exec(
        self, cmd_parts: list[str], cmd_desc: str = "CLI command"
    ) -> CompletedProcess[str]:
        raise NotImplementedError


class IDiskImageService(ICLIService):
    @abstractmethod
    def create_cloud_init_iso(
        self, iso_path: str, paths_to_include: list[str]
    ) -> CompletedProcess[str]:
        raise NotImplementedError


class ILibvirtCLIService(ICLIService):
    @abstractmethod
    def virt_xml_edit(self, vm_name: str, options: list[str]) -> CompletedProcess[str]:
        raise NotImplementedError

    @abstractmethod
    def virt_clone(
        self,
        src_vm_name: str,
        dest_vm_name: str,
        auto_clone: bool = True,
        additional_options: list[str] = None,
    ) -> CompletedProcess[str]:
        raise NotImplementedError
