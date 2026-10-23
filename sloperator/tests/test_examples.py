"""Every program in examples/ must run without errors."""

from pathlib import Path

import pytest

from sloperator.runner import run_source

EXAMPLES = sorted((Path(__file__).resolve().parent.parent / "examples").glob("*.slop"))


@pytest.mark.parametrize("path", EXAMPLES, ids=lambda p: p.stem)
def test_example_runs_cleanly(path: Path):
    result = run_source(path.read_text(encoding="utf-8"))
    assert result.stderr == ""
    assert result.exit_code == 0
    assert result.stdout != ""
