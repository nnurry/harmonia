import subprocess
import logging

from interface.service.cli import IDiskImageService, ILibvirtCLIService, ICLIService

logger = logging.getLogger(__name__)


class CLIService(ICLIService):
    def name(cls):
        return "cli"

    def exec(self, cmd_parts, cmd_desc="CLI command"):
        try:
            result = subprocess.run(
                cmd_parts, check=True, capture_output=True, text=True, shell=False
            )
            if result.stderr:
                logger.warning(f"{cmd_desc} stderr:\n{result.stderr.strip()}")
            else:
                logger.info(f"{cmd_desc} stoud:\n{result.stdout.strip()}")
            logger.info(f"{cmd_desc} completed successfully.")
            return result
        except subprocess.CalledProcessError as e:
            logger.error(f"{cmd_desc} failed with exit code {e.returncode}:")
            logger.error(f"stdout: {e.stdout.strip()}")
            logger.error(f"stderr: {e.stderr.strip()}")
            raise
        except FileNotFoundError:
            logger.error(
                f"Error: Command '{cmd_parts[0]}' not found. Ensure it is installed and in your system's PATH."
            )
            raise
        except Exception as e:
            logger.error(
                f"An unexpected error occurred during {cmd_desc} execution: {e}"
            )
            raise


class DiskImageService(CLIService, IDiskImageService):
    def name(cls):
        return "disk-image-cli"

    def create_cloud_init_iso(self, iso_path, paths_to_include):
        cmd_parts = [
            "mkisofs",
            "-output",
            iso_path,
            "-volid",
            "cidata",
            "-joliet",
            "-r",
        ]
        cmd_parts.extend(paths_to_include)
        return self.exec(cmd_parts, "Cloud-Init ISO creation")


class LibvirtCLIService(CLIService, ILibvirtCLIService):
    def name(cls):
        return "libvirt-cli"

    def virt_xml_edit(self, vm_name, options):
        cmd_parts = ["virt-xml", vm_name, "--edit"]
        cmd_parts.extend(options)
        return self.exec(cmd_parts, f"editing VM '{vm_name}' XML")

    def virt_clone(
        self, src_vm_name, dest_vm_name, auto_clone=True, additional_options=None
    ):
        cmd_parts = [
            "virt-clone",
            "--original",
            src_vm_name,
            "--name",
            dest_vm_name,
        ]
        if auto_clone:
            cmd_parts.append("--auto-clone")
        if additional_options:
            cmd_parts.extend(additional_options)

        return self.exec(cmd_parts, f"cloning VM '{src_vm_name}' to '{dest_vm_name}'")
