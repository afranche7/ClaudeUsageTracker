"""
Claude Usage Tracker — desktop widget for Windows 11
Run: python main.py
"""

import sys
from PyQt6.QtWidgets import QApplication
from src.widget import ClaudeWidget


def main() -> None:
    app = QApplication(sys.argv)
    app.setQuitOnLastWindowClosed(True)

    widget = ClaudeWidget()
    widget.show()

    sys.exit(app.exec())


if __name__ == "__main__":
    main()
