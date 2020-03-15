<?php
include_once '../../../includes/connect_gl.php';
require('../../session/session.php');
include_once '../settings.php';

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  $position = 0;
  $stufenwertung = array(0,0,0,0);
  echo '<ol data-role="listview" data-count-theme="b" data-inset="true">';
  if ($stmt = $mysqli->prepare("SELECT kurz, name, stufe, stamm, sum(points) summe
FROM gruppen, punkte
WHERE gruppen.id = an
GROUP BY an
ORDER BY summe DESC, kurz ASC")) {
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
      $txt .= '<h1>'.$position.'. Platz Gesamtwertung</h1><h2>'.$stufenwertung[$stufe].'. Platz der '.$Stufe[$stufe].'</h2><h3>Mit '.$punkte.' Punkten </h3><h3>"'.$name.'" vom Stamm '.$stamm.' ('.$kurz.')</h3>
      <br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/><br/>';
    }
  }
  echo '</ol>';
}
else {
    echo 'Keine Berechtigung.';
}
