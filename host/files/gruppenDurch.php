<?php
include_once '../../../includes/connect_gl.php';
require('../../session/session.php');
include_once '../settings.php';

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  $file = fopen("../../../siegerehrung/gruppen_durchschnitt.md","w") or die("Einlesen der HTML Datei fehlgeschlagen.");
  $txt = "---
type: slide
slideOptions:
  transition: slide
---

# Georgslauf 2020

Siegerehrung

---
";
  $position = null;
  $stufenwertung = array(0,0,0,0);
  echo '<ol data-role="listview" data-count-theme="b" data-inset="true">';
  if ($stmt = $mysqli->prepare("SELECT kurz, name, stufe, stamm, avg(points) durchschnitt
FROM gruppen, punkte
WHERE gruppen.id = an
GROUP BY an
ORDER BY durchschnitt DESC, kurz ASC")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($kurz, $name, $stufe, $stamm, $punkte);
    while ($stmt->fetch()) {
      if($prev_punkte != $punkte) {
        $position++;
        if($prev_stufe == $stufe)  $stufenwertung[$stufe]--;
      }
      $stufenwertung[$stufe]++;
      echo '<li><a href="#">'.$position.'. Platz ('.$stufenwertung[$stufe].'. der '.$Stufe[$stufe].') "'.$name.'" - '.$stamm.' ('.$kurz.')<span class="ui-li-count">'.$punkte.'</span></a></li>';
      $prev_punkte = $punkte;
      $prev_stufe = $stufe;
      $txt .="
## $position. Platz

### $stufenwertung[$stufe]. Platz der $Stufe[$stufe]
Mit **".round($punkte,2)."** Punkten im Durchschnitt

### *$name* vom Stamm *$stamm*

---
";
    }
  }
  echo '</ol>';
  fwrite($file, $txt);
  fclose($file);
}
