import { themeChange } from 'theme-change'

window.htmx = require('htmx.org');

htmx.config.includeIndicatorStyles = false;
htmx.config.selfRequestsOnly = true;
htmx.config.allowEval = false;

themeChange()
