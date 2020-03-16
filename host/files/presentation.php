<?php
include_once '../../../includes/connect_gl.php';
require('../../session/session.php');
include_once '../settings.php';

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo '<h2>Erstelle Siegertabelle Präsentation</h2>
<button id="copy">Kopieren</button>
<textarea id="input">';
  $md = "---
type: slide
slideOptions:
  transition: slide
---

<!-- .slide: data-background-position=\"bottom\" data-background-size=\"283px\" data-background=\"https://dev.georgslauf.de/res/logo.png\" -->

# Georgslauf 2020

## Siegerehrung

---

# Gruppenwertung

----
";
  $position = 0;
  $prev_stufe = null;
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
# $position. Platz

## $stufenwertung[$stufe]. Platz der $Stufe[$stufe]

## Mit **".round($punkte,2)."** Punkten

### *$name* - *$stamm*

----
".$md_g;
    }
    $md_g .="
<!-- .slide: data-background=\"https://media.giphy.com/media/hqIaXesRGpP44/giphy.gif\" -->

## Herzlichen Glückwunsch 

# Stamm $stamm

## Zum Geogslaufsieg

";
  }
  $md2 = "---

# Postenbewertung

----
";
  $md_p = "";
  $position = 0;
  if ($stmt = $mysqli->prepare("SELECT kurz, name, stufe, stamm, sum(points) as summe FROM posten, punkte WHERE posten.id = an GROUP BY an ORDER BY summe DESC")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($kurz, $name, $stufe, $stamm, $punkte);
    while ($stmt->fetch()) {
      if($prev_punkte != $punkte) $position++;
      $prev_punkte = $punkte;
      if ($position == 1 ) $md_p = "<!-- .slide: data-background=\"gold\" -->".$md_p;
      $md_p = "
# $position. Platz Postenwertung

## Mit **".round($punkte,2)."** Punkten

### *$name* - *$stamm*

----
".$md_p;
    }
    $md_p .="
<!-- .slide: data-background=\"gold\" -->

## Herzlichen Glückwunsch 

# Stamm $stamm

## Zum Postensieg
";
  }

  //writeToFile($md.$md_g.$md2.$md_p);
  echo $md.$md_g.$md2.$md_p;
}
else {
    echo "Keine Berechtigung.";
}
echo "</textarea>";

function writeToFile($text) {
  $file = fopen("../../../siegerehrung/siegerehrung.md","w") or die("Einlesen der MD Datei fehlgeschlagen.");
  fwrite($file, $text);
  fclose($file);
}