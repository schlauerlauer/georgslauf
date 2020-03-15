<?php
include_once '../../../includes/connect_gl.php';
require('../../session/session.php');
include_once '../settings.php';

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  $file = fopen("../../../siegerehrung/siegerehrung.md","w") or die("Einlesen der MD Datei fehlgeschlagen.");
  $md = "---
type: slide
slideOptions:
  transition: slide
---

# Georgslauf 2020

Siegerehrung

---

## Gruppenwertung

----
";
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
      $txt ="
## $position. Platz

### $stufenwertung[$stufe]. Platz der $Stufe[$stufe]

Mit **".round($punkte,2)."** Punkten im Durchschnitt

### *$name* - *$stamm*

----
".$txt;
    }
    $txt .="
<!-- .slide: data-background='https://media.giphy.com/media/hqIaXesRGpP44/giphy.gif' -->

## Herzlichen Glückwunsch 

# Stamm $stamm
";
  }
  fwrite($file, $md.$txt);
  fclose($file);
  echo "Ok.";
}
else {
    echo "Keine Berechtigung.";
}