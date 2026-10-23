from dataclasses import dataclass

from sloperator.errors import SlopRuntimeError
from sloperator.evaluator.values import Value


@dataclass
class Binding:
    value: Value
    mutable: bool
    is_loadbearing: bool = False


class Environment:
    def __init__(self, parent: "Environment | None" = None) -> None:
        self.bindings: dict[str, Binding] = {}
        self.parent = parent

    def declare(
        self, name: str, value: Value, mutable: bool, line: int, is_loadbearing: bool = False
    ) -> None:
        if name in self.bindings:
            raise SlopRuntimeError(line, f"'{name}' is already declared in this scope")
        self.bindings[name] = Binding(value, mutable, is_loadbearing)

    def get(self, name: str, line: int) -> Value:
        return self.lookup(name, line).value

    def assign(self, name: str, value: Value, line: int) -> None:
        binding = self.lookup(name, line)
        if not binding.mutable:
            if binding.is_loadbearing:
                raise SlopRuntimeError(line, f"cannot assign to loadbearing '{name}'")
            raise SlopRuntimeError(line, f"cannot assign to immutable binding '{name}'")
        binding.value = value

    def lookup(self, name: str, line: int) -> Binding:
        env: Environment | None = self
        while env is not None:
            if name in env.bindings:
                return env.bindings[name]
            env = env.parent
        raise SlopRuntimeError(line, f"undefined variable '{name}'")
