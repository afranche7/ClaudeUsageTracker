import sys
import time
from datetime import datetime

from PyQt6.QtCore import Qt, QTimer, QPoint
from PyQt6.QtGui import (
    QColor, QPainter, QPainterPath, QPen, QFont, QFontDatabase,
    QCursor, QAction,
)
from PyQt6.QtWidgets import QWidget, QApplication, QMenu

from .tracker import (
    get_current_session_id,
    aggregate_session,
    aggregate_weekly,
    format_tokens,
    format_cost,
)

# ── Claude brand palette ────────────────────────────────────────────────────
BG_COLOR = QColor(15, 10, 28, 224)          # deep navy-purple, ~88% opaque
BORDER_COLOR = QColor(217, 119, 87, 60)     # faint coral glow
ACCENT = QColor(217, 119, 87)               # #d97757  warm coral-orange
TEXT_PRIMARY = QColor(240, 235, 228)        # #f0ebe4  warm cream
TEXT_MUTED = QColor(123, 116, 144)          # #7b7490  soft purple-gray
DIVIDER = QColor(217, 119, 87, 35)

CORNER_RADIUS = 14
WIDGET_W = 320
WIDGET_H = 175
REFRESH_MS = 30_000   # 30 seconds


def _try_acrylic(hwnd: int) -> bool:
    """
    Attempt Windows 11 Acrylic backdrop. Tries two approaches:
    1. DwmSetWindowAttribute (Windows 11 22H2+)
    2. SetWindowCompositionAttribute (Windows 10/11 legacy fallback)
    Silent no-op on non-Windows.
    """
    import sys
    if sys.platform != "win32":
        return False

    # Approach 1: Windows 11 22H2+ DWM system backdrop
    try:
        import ctypes
        import ctypes.wintypes
        dwmapi = ctypes.windll.dwmapi  # type: ignore[attr-defined]
        DWMWA_SYSTEMBACKDROP_TYPE = 38
        DWMWCP_ACRYLIC = 3
        value = ctypes.c_int(DWMWCP_ACRYLIC)
        result = dwmapi.DwmSetWindowAttribute(
            ctypes.wintypes.HWND(hwnd),
            DWMWA_SYSTEMBACKDROP_TYPE,
            ctypes.byref(value),
            ctypes.sizeof(value),
        )
        if result == 0:
            return True
    except Exception:
        pass

    # Approach 2: SetWindowCompositionAttribute (Windows 10/11 legacy)
    try:
        import ctypes

        class ACCENTPOLICY(ctypes.Structure):
            _fields_ = [
                ("AccentState", ctypes.c_int),
                ("AccentFlags", ctypes.c_int),
                ("GradientColor", ctypes.c_uint),
                ("AnimationId", ctypes.c_int),
            ]

        class WINCOMPATTRDATA(ctypes.Structure):
            _fields_ = [
                ("Attribute", ctypes.c_int),
                ("Data", ctypes.c_void_p),
                ("SizeOfData", ctypes.c_size_t),
            ]

        accent = ACCENTPOLICY()
        accent.AccentState = 4           # ACCENT_ENABLE_ACRYLICBLURBEHIND
        accent.AccentFlags = 2
        accent.GradientColor = 0xD00F0A1C  # #0F0A1C at ~82% opacity (AABBGGRR)

        data = WINCOMPATTRDATA()
        data.Attribute = 19              # WCA_ACCENT_POLICY
        data.SizeOfData = ctypes.sizeof(accent)
        data.Data = ctypes.cast(ctypes.pointer(accent), ctypes.c_void_p)

        user32 = ctypes.windll.user32  # type: ignore[attr-defined]
        user32.SetWindowCompositionAttribute(hwnd, ctypes.byref(data))
        return True
    except Exception:
        pass

    return False


class ClaudeWidget(QWidget):
    def __init__(self) -> None:
        super().__init__()

        self._drag_pos: QPoint | None = None
        self._last_refresh: float = 0.0
        self._session_data: dict = {}
        self._weekly_data: dict = {}
        self._session_id: str | None = None

        self._setup_window()
        self._setup_fonts()
        self._position_bottom_right()

        # Initial data load
        self._refresh()

        # Periodic timer
        self._timer = QTimer(self)
        self._timer.timeout.connect(self._refresh)
        self._timer.start(REFRESH_MS)

    # ── Window setup ────────────────────────────────────────────────────────

    def _setup_window(self) -> None:
        self.setWindowFlags(
            Qt.WindowType.FramelessWindowHint
            | Qt.WindowType.WindowStaysOnTopHint
            | Qt.WindowType.Tool              # no taskbar entry
        )
        self.setAttribute(Qt.WidgetAttribute.WA_TranslucentBackground)
        self.setAttribute(Qt.WidgetAttribute.WA_NoSystemBackground)
        self.setFixedSize(WIDGET_W, WIDGET_H)

    def showEvent(self, event):
        super().showEvent(event)
        hwnd = int(self.winId())
        _try_acrylic(hwnd)

    def _position_bottom_right(self) -> None:
        screen = QApplication.primaryScreen()
        if screen:
            geo = screen.availableGeometry()
            margin = 20
            self.move(
                geo.right() - WIDGET_W - margin,
                geo.bottom() - WIDGET_H - margin,
            )

    def _setup_fonts(self) -> None:
        self._font_header = QFont("Segoe UI", 9, QFont.Weight.DemiBold)
        self._font_label = QFont("Segoe UI", 7, QFont.Weight.Normal)
        self._font_value = QFont("Segoe UI", 16, QFont.Weight.Bold)
        self._font_cost = QFont("Segoe UI", 9, QFont.Weight.Normal)
        self._font_footer = QFont("Segoe UI", 7, QFont.Weight.Normal)

    # ── Painting ────────────────────────────────────────────────────────────

    def paintEvent(self, event) -> None:
        p = QPainter(self)
        p.setRenderHint(QPainter.RenderHint.Antialiasing)

        # Background rounded rect
        path = QPainterPath()
        path.addRoundedRect(0, 0, WIDGET_W, WIDGET_H, CORNER_RADIUS, CORNER_RADIUS)
        p.fillPath(path, BG_COLOR)

        # Outer border
        p.setPen(QPen(BORDER_COLOR, 1.0))
        p.drawRoundedRect(1, 1, WIDGET_W - 2, WIDGET_H - 2, CORNER_RADIUS, CORNER_RADIUS)

        # ── Header row ────────────────────────────────────────────────────
        pad = 16
        y = 14

        # Orange diamond dot
        p.setBrush(ACCENT)
        p.setPen(Qt.PenStyle.NoPen)
        dot_size = 8
        dot_x = pad
        dot_y = y - dot_size // 2 + 1
        p.save()
        p.translate(dot_x + dot_size / 2, dot_y + dot_size / 2)
        p.rotate(45)
        p.drawRect(-dot_size // 2, -dot_size // 2, dot_size, dot_size)
        p.restore()

        # "Claude Usage" text
        p.setFont(self._font_header)
        p.setPen(TEXT_PRIMARY)
        p.drawText(pad + dot_size + 7, y + 4, "Claude Usage")

        # Close button "×"
        p.setFont(QFont("Segoe UI", 10))
        p.setPen(TEXT_MUTED)
        p.drawText(WIDGET_W - pad - 10, y + 4, "×")

        # Divider line under header
        div_y = y + 14
        p.setPen(QPen(DIVIDER, 1))
        p.drawLine(pad, div_y, WIDGET_W - pad, div_y)

        # ── Two-column content ─────────────────────────────────────────────
        col_left = pad
        col_right = WIDGET_W // 2 + 8
        content_y = div_y + 16

        self._draw_column(p, col_left, content_y, "SESSION", self._session_data)
        self._draw_column(p, col_right, content_y, "WEEKLY", self._weekly_data)

        # Vertical divider between columns
        mid_x = WIDGET_W // 2
        p.setPen(QPen(DIVIDER, 1))
        p.drawLine(mid_x, div_y + 4, mid_x, div_y + 90)

        # ── Footer ────────────────────────────────────────────────────────
        if self._last_refresh:
            elapsed = int(time.time() - self._last_refresh)
            if elapsed < 60:
                age = f"updated {elapsed}s ago"
            else:
                age = f"updated {elapsed // 60}m ago"
        else:
            age = "loading…"

        p.setFont(self._font_footer)
        p.setPen(TEXT_MUTED)
        p.drawText(pad, WIDGET_H - 10, age)

    def _draw_column(self, p: QPainter, x: int, y: int, label: str, data: dict) -> None:
        # Column label
        p.setFont(self._font_label)
        p.setPen(TEXT_MUTED)
        p.drawText(x, y, label)

        tokens = format_tokens(data.get("total_tokens", 0)) if data else "—"
        cost = format_cost(data.get("cost_usd", 0.0)) if data else "—"

        # Token value (large)
        p.setFont(self._font_value)
        p.setPen(TEXT_PRIMARY)
        p.drawText(x, y + 32, tokens)

        # "tokens" sub-label
        p.setFont(self._font_label)
        p.setPen(TEXT_MUTED)
        p.drawText(x, y + 44, "tokens")

        # Cost value
        p.setFont(self._font_cost)
        p.setPen(ACCENT)
        p.drawText(x, y + 64, cost)

    # ── Data refresh ────────────────────────────────────────────────────────

    def _refresh(self) -> None:
        try:
            self._session_id = get_current_session_id()
            self._session_data = aggregate_session(self._session_id)
            self._weekly_data = aggregate_weekly()
            self._last_refresh = time.time()
        except Exception:
            pass
        self.update()

    # ── Interaction ─────────────────────────────────────────────────────────

    def mousePressEvent(self, event) -> None:
        if event.button() == Qt.MouseButton.LeftButton:
            # Check if close button clicked (top-right ~20x20 area)
            if event.position().x() > WIDGET_W - 36 and event.position().y() < 30:
                QApplication.quit()
                return
            self._drag_pos = event.globalPosition().toPoint() - self.frameGeometry().topLeft()
            event.accept()

    def mouseMoveEvent(self, event) -> None:
        if self._drag_pos and event.buttons() == Qt.MouseButton.LeftButton:
            self.move(event.globalPosition().toPoint() - self._drag_pos)
            event.accept()

    def mouseReleaseEvent(self, event) -> None:
        self._drag_pos = None

    def contextMenuEvent(self, event) -> None:
        menu = QMenu(self)
        menu.setStyleSheet("""
            QMenu {
                background-color: #120d20;
                color: #f0ebe4;
                border: 1px solid rgba(217,119,87,0.3);
                border-radius: 6px;
                padding: 4px;
                font-family: 'Segoe UI';
                font-size: 9pt;
            }
            QMenu::item {
                padding: 5px 18px;
                border-radius: 4px;
            }
            QMenu::item:selected {
                background-color: rgba(217,119,87,0.18);
            }
            QMenu::separator {
                height: 1px;
                background: rgba(217,119,87,0.2);
                margin: 3px 8px;
            }
        """)

        refresh_action = QAction("Refresh now", self)
        refresh_action.triggered.connect(self._refresh)

        reset_action = QAction("Reset session", self)
        reset_action.triggered.connect(self._reset_session)

        menu.addAction(refresh_action)
        menu.addAction(reset_action)
        menu.addSeparator()

        exit_action = QAction("Exit", self)
        exit_action.triggered.connect(QApplication.quit)
        menu.addAction(exit_action)

        menu.exec(event.globalPos())

    def _reset_session(self) -> None:
        """Clear session data (useful if you want to start a fresh tracking window)."""
        self._session_data = {}
        self._last_refresh = time.time()
        self.update()
