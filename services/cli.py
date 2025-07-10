import subprocess
import logging

logger = logging.getLogger(__name__)


class CLI:
    @staticmethod
    def run_command(command_parts: list[str], description: str = "CLI command"):
        try:
            result = subprocess.run(
                command_parts, check=True, capture_output=True, text=True, shell=False
            )
            if result.stderr:
                logger.warning(f"{description} stderr:\n{result.stderr.strip()}")
            else:
                logger.info(f"{description} stoud:\n{result.stdout.strip()}")
            logger.info(f"{description} completed successfully.")
            return result
        except subprocess.CalledProcessError as e:
            logger.error(f"{description} failed with exit code {e.returncode}:")
            logger.error(f"stdout: {e.stdout.strip()}")
            logger.error(f"stderr: {e.stderr.strip()}")
            raise
        except FileNotFoundError:
            logger.error(
                f"Error: Command '{command_parts[0]}' not found. Ensure it is installed and in your system's PATH."
            )
            raise
        except Exception as e:
            logger.error(
                f"An unexpected error occurred during {description} execution: {e}"
            )
            raise


class DiskImageService(CLI):
    @staticmethod
    def create_cloud_init_iso(iso_path: str, paths_to_include: list[str]) -> str:
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
        return CLI.run_command(cmd_parts, "Cloud-Init ISO creation")


class LibvirtCLIService(CLI):
    @staticmethod
    def virt_xml_edit(vm_name: str, options: list[str]):
        cmd_parts = ["virt-xml", vm_name, "--edit"]
        cmd_parts.extend(options)
        return CLI.run_command(cmd_parts, f"editing VM '{vm_name}' XML")

    @staticmethod
    def virt_clone(
        src_vm_name: str,
        dest_vm_name: str,
        auto_clone: bool = True,
        additional_options: list[str] = None,
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
        return CLI.run_command(
            cmd_parts, f"cloning VM '{src_vm_name}' to '{dest_vm_name}'"
        )
