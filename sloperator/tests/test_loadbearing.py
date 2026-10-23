from sloperator.lexer import TokenType, tokenize
from sloperator.parser import parse
from sloperator.parser.ast import VarDecl
from sloperator.runner import run_source


def run(source: str) -> str:
    result = run_source(source)
    assert result.stderr == "", result.stderr
    return result.stdout


def fails(source: str) -> str:
    result = run_source(source)
    assert result.exit_code == 1
    return result.stderr.rstrip("\n")


def test_loadbearing_is_a_keyword():
    (token,) = [t for t in tokenize("loadbearing") if t.type is not TokenType.EOF]
    assert token.type is TokenType.LOADBEARING


def test_loadbearing_declaration_is_immutable():
    (declaration,) = parse(tokenize("loadbearing LIMIT = 100✨"))
    assert isinstance(declaration, VarDecl) and not declaration.mutable


def test_loadbearing_can_be_read():
    assert run("loadbearing LIMIT = 100✨\nhallucinate(LIMIT)✨") == "100\n"


def test_slop_can_still_be_reassigned():
    source = (
        "loadbearing LIMIT = 100✨\nslop score = 10✨\nscore = score + 1✨\nhallucinate(score)✨"
    )
    assert run(source) == "11\n"


def test_loadbearing_reassignment_fails():
    assert fails("loadbearing LIMIT = 100✨\nLIMIT = 200✨") == (
        "SlopLang runtime error at line 2: cannot assign to loadbearing 'LIMIT'"
    )


def test_loadbearing_reassignment_fails_inside_block():
    source = "loadbearing LIMIT = 100✨\n{\n    LIMIT = 200✨\n}"
    assert fails(source) == (
        "SlopLang runtime error at line 3: cannot assign to loadbearing 'LIMIT'"
    )


def test_loadbearing_assignment_hits_nearest_binding():
    source = "slop x = 1✨\n{\n    loadbearing x = 2✨\n    x = 3✨\n}"
    assert fails(source) == ("SlopLang runtime error at line 4: cannot assign to loadbearing 'x'")


def test_loadbearing_cannot_be_redeclared_in_same_scope():
    assert fails("loadbearing x = 1✨\nloadbearing x = 2✨") == (
        "SlopLang runtime error at line 2: 'x' is already declared in this scope"
    )
    assert fails("slop x = 1✨\nloadbearing x = 2✨") == (
        "SlopLang runtime error at line 2: 'x' is already declared in this scope"
    )
    assert fails("loadbearing x = 1✨\nslop x = 2✨") == (
        "SlopLang runtime error at line 2: 'x' is already declared in this scope"
    )


def test_loadbearing_may_shadow_outer_scope():
    source = (
        "loadbearing x = 1✨\n{\n    loadbearing x = 2✨\n    hallucinate(x)✨\n}\nhallucinate(x)✨"
    )
    assert run(source) == "2\n1\n"


def test_slop_may_shadow_outer_loadbearing():
    source = (
        "loadbearing x = 1✨\n"
        "{\n"
        "    slop x = 2✨\n"
        "    x = 3✨\n"
        "    hallucinate(x)✨\n"
        "}\n"
        "hallucinate(x)✨"
    )
    assert run(source) == "3\n1\n"


def test_loadbearing_requires_initializer():
    result = run_source("loadbearing x✨")
    assert result.exit_code == 1
    assert result.stderr == (
        "SlopLang syntax error at line 1: "
        "expected '=' after variable name, an initializer is required, found '✨'\n"
    )
