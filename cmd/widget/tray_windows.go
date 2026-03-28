//go:build windows

package main

import (
	"fmt"
	"os"

	"github.com/afranche7/claudeusagetracker/internal/tracker"
	"github.com/afranche7/claudeusagetracker/internal/ui"
	"github.com/getlantern/systray"
)

func initSystemTray() {
	systray.Run(onTrayReady, onTrayExit)
}

func onTrayReady() {
	systray.SetIcon(generateMinimalIcon())
	systray.SetTitle("Claude Usage")
	systray.SetTooltip("Claude Code Usage Tracker")

	mRefresh := systray.AddMenuItem("Refresh Now", "Refresh usage data")
	systray.AddSeparator()

	mPlan := systray.AddMenuItem("Plan: "+tracker.GetPlanLimits().Label, "Current plan")
	mPlan.Disable()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Exit the tracker")

	go func() {
		for {
			select {
			case <-mRefresh.ClickedCh:
				data := ui.BuildWidgetData()
				plan := tracker.GetPlanLimits()
				systray.SetTooltip(fmt.Sprintf(
					"Today: $%s / $%.0f | Week: $%s / $%.0f",
					data.TodayCost, plan.DailyUSD,
					data.WeeklyCost, plan.WeeklyUSD,
				))
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func onTrayExit() {}

// generateMinimalIcon creates a simple 16x16 ICO with Claude's tan accent color
func generateMinimalIcon() []byte {
	width, height := 16, 16
	bmpHeaderSize := 40
	pixelDataSize := width * height * 4
	maskSize := ((width + 31) / 32) * 4 * height
	imageSize := bmpHeaderSize + pixelDataSize + maskSize
	totalSize := 6 + 16 + imageSize

	ico := make([]byte, totalSize)

	// ICO header
	ico[2] = 1 // type: ICO
	ico[4] = 1 // count: 1

	// Directory entry
	ico[6] = byte(width)
	ico[7] = byte(height)
	ico[10] = 1  // color planes
	ico[12] = 32 // bpp
	ico[14] = byte(imageSize)
	ico[15] = byte(imageSize >> 8)
	offset := 22
	ico[18] = byte(offset)

	// BITMAPINFOHEADER
	o := offset
	ico[o] = 40
	ico[o+4] = byte(width)
	ico[o+8] = byte(height * 2)
	ico[o+12] = 1  // planes
	ico[o+14] = 32 // bpp
	o += bmpHeaderSize

	// Pixel data — Claude "C" logo in tan (#D4A27C)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := o + ((height-1-y)*width+x)*4
			if idx+3 >= len(ico) {
				continue
			}
			cx, cy := float64(x)-7.5, float64(y)-7.5
			dist := cx*cx + cy*cy
			inRing := dist < 49 && dist > 20
			isC := inRing && !(cx > 1 && cy > -3 && cy < 3)
			if isC {
				ico[idx] = 124   // B
				ico[idx+1] = 162 // G
				ico[idx+2] = 212 // R
				ico[idx+3] = 255 // A
			}
		}
	}

	return ico
}
