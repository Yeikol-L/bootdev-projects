from collections.abc import Callable

from sloperator.parser.ast import FunctionDecl


class SlopFunction:
    def __init__(self, declaration: FunctionDecl) -> None:
        self.declaration = declaration

    @property
    def name(self) -> str:
        return self.declaration.name.lexeme


class Builtin:
    def __init__(self, name: str, call: Callable[[list["Value"], int], "Value"]) -> None:
        self.name = name
        self.call = call


type Value = int | str | bool | None | SlopFunction | Builtin


def type_name(value: Value) -> str:
    if isinstance(value, bool):
        return "boolean"
    if isinstance(value, int):
        return "number"
    if isinstance(value, str):
        return "string"
    if value is None:
        return "nil"
    return "clanker"


def stringify(value: Value) -> str:
    if value is True:
        return "✅"
    if value is False:
        return "❌"
    if value is None:
        return "nil"
    if isinstance(value, SlopFunction | Builtin):
        return f"<clanker {value.name}>"
    return str(value)


def values_equal(left: Value, right: Value) -> bool:
    if isinstance(left, SlopFunction | Builtin) or isinstance(right, SlopFunction | Builtin):
        return left is right

    return type(left) is type(right) and left == right
