"""
Claude Usage Tracker — desktop widget for Windows 11
Run: python main.py
"""

import sys

# Must run before QApplication on Windows to get crisp rendering on HiDPI screens
if sys.platform == "win32":
    try:
        import ctypes
        ctypes.windll.shcore.SetProcessDpiAwareness(2)  # PROCESS_PER_MONITOR_DPI_AWARE
    except Exception:
        pass

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
