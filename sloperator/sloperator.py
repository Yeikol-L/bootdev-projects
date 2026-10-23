import sys
from pathlib import Path

from sloperator.runner import run_program, validate_program


def main(argv: list[str]) -> int:
    args = argv[1:]
    check = "--check" in args
    files = [arg for arg in args if arg != "--check"]
    if any(arg.startswith("-") for arg in files) or len(files) != 1:
        print("usage: uv run sloperator.py [--check] <file.slop>", file=sys.stderr)
        return 2

    path = Path(files[0])
    try:
        source = path.read_text(encoding="utf-8")
    except OSError as error:
        print(f"could not read '{path}': {error.strerror}", file=sys.stderr)
        return 2

    if check:
        error = validate_program(source)
    else:
        error = run_program(source, sys.stdout)
    if error is not None:
        sys.stdout.flush()
        print(error, file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv))
