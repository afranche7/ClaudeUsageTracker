package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/afranche7/claudeusagetracker/internal/ui"
)

var (
	flagPort     = flag.Int("port", 0, "HTTP port for browser mode (0 = auto)")
	flagInterval = flag.Int("interval", 30, "Refresh interval in seconds")
	flagBrowser  = flag.Bool("browser", false, "Force browser mode instead of native widget")
)

func main() {
	flag.Parse()

	// Try native webview first (Windows), fall back to browser mode
	if !*flagBrowser {
		if launchNativeWidget(*flagInterval) {
			return
		}
		log.Println("Native widget not available, falling back to browser mode")
	}

	launchBrowserMode(*flagPort, *flagInterval)
}

// launchBrowserMode starts an HTTP server and opens the widget in a browser tab
func launchBrowserMode(port, intervalSec int) {
	if port == 0 {
		port = 17429 // Default port
	}

	var mu sync.RWMutex
	var currentHTML string

	// Initial render
	data := ui.BuildWidgetData()
	currentHTML = ui.RenderHTML(data)

	// Background refresh
	go func() {
		ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			newData := ui.BuildWidgetData()
			html := ui.RenderHTML(newData)
			mu.Lock()
			currentHTML = html
			mu.Unlock()
		}
	}()

	// Widget page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		html := currentHTML
		mu.RUnlock()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, wrapWithAutoRefresh(html, intervalSec))
	})

	// JSON API endpoint for programmatic access
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		data := ui.BuildWidgetData()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	})

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("Shutting down...")
		os.Exit(0)
	}()

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	log.Printf("Claude Usage Tracker running at http://%s", addr)
	log.Printf("Refreshing every %ds. Press Ctrl+C to quit.", intervalSec)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// wrapWithAutoRefresh adds a meta refresh and auto-reload script to the HTML
func wrapWithAutoRefresh(html string, intervalSec int) string {
	script := fmt.Sprintf(`<script>
setTimeout(function(){
  fetch(window.location.href).then(r=>r.text()).then(h=>{
    let parser = new DOMParser();
    let doc = parser.parseFromString(h, 'text/html');
    let newWidget = doc.querySelector('.widget');
    if(newWidget) {
      document.querySelector('.widget').innerHTML = newWidget.innerHTML;
    }
  });
  setInterval(function(){
    fetch(window.location.href).then(r=>r.text()).then(h=>{
      let parser = new DOMParser();
      let doc = parser.parseFromString(h, 'text/html');
      let newWidget = doc.querySelector('.widget');
      if(newWidget) {
        document.querySelector('.widget').innerHTML = newWidget.innerHTML;
      }
    });
  }, %d);
}, %d);
</script>`, intervalSec*1000, intervalSec*1000)

	// Insert script before </body>
	return html + script
}
