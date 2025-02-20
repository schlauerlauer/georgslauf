// HTMX
import htmx from "htmx.org";

htmx.config.includeIndicatorStyles = false;
htmx.config.selfRequestsOnly = true;
htmx.config.allowEval = false;
htmx.config.allowScriptTags = false;

// HYPERSCRIPT
// import _hyperscript from "hyperscript.org";
// _hyperscript.browserInit();

// THEMES
import { themeChange } from 'theme-change'
themeChange()
