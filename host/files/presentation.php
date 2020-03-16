<?php
include_once '../../../includes/connect_gl.php';
require('../../session/session.php');
include_once '../settings.php';

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo '<h2>Erstelle Siegertabelle Präsentation</h2>
<button id="copy">Kopieren</button>
<textarea id="input">';
  echo startSlide();//.gruppenSlide();
  echo "</textarea>";
}
else {
  echo "Keine Berechtigung.";
}

function writeToFile($write) {
  $file = fopen("../../../siegerehrung/siegerehrung.md","w") or die("Einlesen der MD Datei fehlgeschlagen.");
  fwrite($file, $write);
  fclose($file);
}
function startSlide() {
  return "---
  type: slide
  slideOptions:
    transition: slide
  ---
  
  # Georgslauf 2020
  
  Siegerehrung
  
  ---

";
}

function postenSlide() {
  return "---

## Postenbewertung

----
";
}