from abc import ABCMeta, abstractmethod


class IService(metaclass=ABCMeta):
    @abstractmethod
    def name(cls) -> str:
        raise NotImplementedError
