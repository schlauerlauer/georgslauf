<?php
include_once '../../../includes/connect_gl.php';
require('../../session/session.php');
include_once '../settings.php';

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  $file = fopen("../../../siegerehrung/gruppensieger.html","w") or die("Einlesen der HTML Datei fehlgeschlagen.");
  $txt = "";
  $position = null;
  echo '<ol data-role="listview" data-count-theme="b" data-inset="true">';
  if ($stmt = $mysqli->prepare("SELECT kurz, name, stufe, stamm, avg(points) as durchschnitt FROM posten, punkte WHERE posten.id = an GROUP BY an ORDER BY durchschnitt DESC")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($kurz, $name, $stufe, $stamm, $punkte);
    while ($stmt->fetch()) {
      if($prev_punkte != $punkte) $position++;
      echo '<li><a href="#">'.$position.'. Platz - '.$stufe.' '.$kurz.' - "'.$name.'" - '.$stamm.'<span class="ui-li-count">'.$punkte.'</span></a></li>';
      $prev_punkte = $punkte;
      $txt .= '<h1>'.$position.'. Platz Postenwertung</h1><h3>Mit '.$punkte.' Punkten im Durchschnitt</h3><h3>"'.$name.'" vom Stamm '.$stamm.' ('.$stufe.' '.$kurz.')</h3>
      <br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/>';
    }
  }
  echo '</ol>';
  fwrite($file, $txt);
  fclose($file);
}
