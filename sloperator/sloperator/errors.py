class SlopError(Exception):
    kind = "error"

    def __init__(self, line: int, message: str) -> None:
        super().__init__(message)
        self.line = line
        self.message = message

    def __str__(self) -> str:
        return f"SlopLang {self.kind} error at line {self.line}: {self.message}"


class LexicalError(SlopError):
    kind = "lexical"


class SlopSyntaxError(SlopError):
    kind = "syntax"


class SlopRuntimeError(SlopError):
    kind = "runtime"
