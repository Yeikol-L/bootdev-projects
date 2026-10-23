from sloperator.runner import run_source


def test_multiplication_before_addition():
    result = run_source("hallucinate(2 + 3 * 4)✨")
    assert result.stderr == ""
    assert result.stdout == "14\n"
