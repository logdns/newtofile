package appui

import "github.com/egoist/mygo/ui"

func icon(paths string) *ui.SVG {
	return ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">` + paths + `</svg>`))
}

var (
	codeIcon         = icon(`<path d="M8 5 2 12l6 7M16 5l6 7-6 7M14 3l-4 18"/>`)
	tableIcon        = icon(`<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M3 10h18M3 15h18M10 4v16"/>`)
	presentationIcon = icon(`<rect x="3" y="3" width="18" height="13" rx="2"/><path d="M12 16v5M8 21l4-3 4 3M7 8h10M7 12h5"/>`)
	documentIcon     = icon(`<path d="M14 3H6a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><path d="M14 3v6h6M8 13h8M8 17h5"/>`)
	plusIcon         = icon(`<path d="M12 5v14M5 12h14"/>`)
	folderIcon       = icon(`<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>`)
	gridIcon         = icon(`<rect x="3" y="3" width="7" height="7" rx="1.5"/><rect x="14" y="3" width="7" height="7" rx="1.5"/><rect x="3" y="14" width="7" height="7" rx="1.5"/><rect x="14" y="14" width="7" height="7" rx="1.5"/>`)
	officeIcon       = icon(`<rect x="4" y="3" width="16" height="18" rx="2"/><path d="M8 7h2M14 7h2M8 11h2M14 11h2M10 21v-6h4v6"/>`)
	desktopIcon      = icon(`<rect x="3" y="3" width="18" height="13" rx="2"/><path d="M8 21h8M12 16v5"/>`)
	downloadIcon     = icon(`<path d="M12 3v12M7 10l5 5 5-5M4 16v4h16v-4"/>`)
	settingsIcon     = icon(`<path d="M4 7h16M4 17h16"/><circle cx="8" cy="7" r="3"/><circle cx="16" cy="17" r="3"/>`)
	refreshIcon      = icon(`<path d="M20 7v5h-5M4 17v-5h5M6 6a8 8 0 0 1 13 2M18 18a8 8 0 0 1-13-2"/>`)
	checkIcon        = icon(`<path d="m5 12 4 4L19 6"/>`)
	arrowIcon        = icon(`<path d="M7 17 17 7M7 7h10v10"/>`)
	lockIcon         = icon(`<rect x="5" y="10" width="14" height="11" rx="2"/><path d="M8 10V7a4 4 0 0 1 8 0v3M12 14v3"/>`)
)

func formatIcon(kind string) *ui.SVG {
	switch kind {
	case "json", "html":
		return codeIcon
	case "csv", "xlsx":
		return tableIcon
	case "pptx":
		return presentationIcon
	default:
		return documentIcon
	}
}
