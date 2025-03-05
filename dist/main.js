// HTMX
import htmx from "htmx.org";

htmx.config.includeIndicatorStyles = false;
htmx.config.selfRequestsOnly = true;
htmx.config.allowEval = false;
htmx.config.allowScriptTags = false;

htmx.on("htmx:beforeOnLoad", function (event) {
	if (!event.detail) {
		return
	}
	if (event.detail.xhr.status >= 400) {
		event.detail.shouldSwap = true;
		event.detail.isError = false;
	}
});

// HYPERSCRIPT
import _hyperscript from "hyperscript.org";
_hyperscript.browserInit();

// THEMES
import { themeChange } from 'theme-change'
themeChange()
