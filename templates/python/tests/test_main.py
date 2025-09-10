# SPDX-FileCopyrightText: 2025-present Your Name <your.email@example.com>
#
# SPDX-License-Identifier: MIT

"""Tests for the main module."""

import pytest
from my_python_project.main import main


def test_main_runs_without_error(capsys):
    """Test that main function runs without error and produces expected output."""
    main()
    captured = capsys.readouterr()
    assert "Welcome to your new Python project!" in captured.out
    assert "TGS workflow" in captured.out