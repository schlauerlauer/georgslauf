<?php
include_once '../../../includes/connect_gl.php';
require('../../session/session.php');
include_once '../settings.php';

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo '<h2>Erstelle Siegertabelle Präsentation</h2>
<button id="copy">Kopieren</button>
<textarea id="input">';
  echo startSlide().gruppenSlide();
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
function gruppenSlide() {
    $position = 0;
    $stufenwertung = array(0,0,0,0);
    if ($stmt = $mysqli->prepare("SELECT kurz, name, stufe, stamm, sum(points) summe FROM gruppen, punkte WHERE gruppen.id = an GROUP BY an ORDER BY summe DESC, kurz ASC")) {
      $stmt->execute();
      $stmt->store_result();
      $stmt->bind_result($kurz, $name, $stufe, $stamm, $punkte);
      while ($stmt->fetch()) {
        if($prev_punkte != $punkte) {
          $position++;
          if($prev_stufe == $stufe)  $stufenwertung[$stufe]--;
        }
        $stufenwertung[$stufe]++;
        $prev_punkte = $punkte;
        $prev_stufe = $stufe;
        $md_g ="
  ## $position. Platz
  
  ### $stufenwertung[$stufe]. Platz der $Stufe[$stufe]
  
  Mit **".round($punkte,2)."** Punkten im Durchschnitt
  
  ### *$name* - *$stamm*
  
  ----
  ".$md_g;
      }
      $md_g .="
  <!-- .slide: data-background=\"https://media.giphy.com/media/hqIaXesRGpP44/giphy.gif\" -->
  
  ## Herzlichen Glückwunsch 
  
  # Stamm $stamm
  
  ";
  $md_g = "## Gruppenwertung
  
  ----".$md_g;
  return $md_g;
}
function postenSlide() {
  $md2 = "---

## Postenbewertung

----
";
  return $md2;
}