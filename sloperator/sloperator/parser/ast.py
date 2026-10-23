from dataclasses import dataclass

from sloperator.lexer.tokens import Token


@dataclass(frozen=True)
class Expr:
    @property
    def line(self) -> int:
        raise NotImplementedError


@dataclass(frozen=True)
class Literal(Expr):
    value: int | str | bool
    token: Token

    @property
    def line(self) -> int:
        return self.token.line


@dataclass(frozen=True)
class Variable(Expr):
    name: Token

    @property
    def line(self) -> int:
        return self.name.line


@dataclass(frozen=True)
class Unary(Expr):
    operator: Token
    operand: Expr

    @property
    def line(self) -> int:
        return self.operator.line


@dataclass(frozen=True)
class Binary(Expr):
    left: Expr
    operator: Token
    right: Expr

    @property
    def line(self) -> int:
        return self.operator.line


@dataclass(frozen=True)
class Logical(Expr):
    left: Expr
    operator: Token
    right: Expr

    @property
    def line(self) -> int:
        return self.operator.line


@dataclass(frozen=True)
class Call(Expr):
    callee: Expr
    paren: Token
    arguments: tuple[Expr, ...]

    @property
    def line(self) -> int:
        return self.paren.line


@dataclass(frozen=True)
class Grouping(Expr):
    expression: Expr

    @property
    def line(self) -> int:
        return self.expression.line


@dataclass(frozen=True)
class Stmt:
    pass


@dataclass(frozen=True)
class Block(Stmt):
    statements: tuple[Stmt, ...]


@dataclass(frozen=True)
class VarDecl(Stmt):
    name: Token
    initializer: Expr
    mutable: bool


@dataclass(frozen=True)
class Assign(Stmt):
    name: Token
    operator: Token
    value: Expr


@dataclass(frozen=True)
class ExprStmt(Stmt):
    expression: Expr


@dataclass(frozen=True)
class If(Stmt):
    condition: Expr
    then_branch: Block
    else_branch: Block | None


@dataclass(frozen=True)
class While(Stmt):
    condition: Expr
    body: Block


@dataclass(frozen=True)
class Ship(Stmt):
    keyword: Token
    value: Expr


@dataclass(frozen=True)
class FunctionDecl(Stmt):
    name: Token
    params: tuple[Token, ...]
    body: Block


Program = list[Stmt]
