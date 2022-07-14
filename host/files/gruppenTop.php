<?php
require('../../session/session.php');
include_once '../settings.php';

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo "<h2>Gruppen Sieger</h2>";
  $position = 0;
  $stufenwertung = array(0,0,0,0);
  $stufenpunkte = array(0,0,0,0);
  echo '<ul data-role="listview" data-count-theme="b" data-inset="true">';
  if ($stmt = $mysqli->prepare("SELECT kurz, name, stufe, stamm, sum(points) summe
FROM gruppen, punkte
WHERE gruppen.id = an
GROUP BY an
ORDER BY summe DESC, stufe ASC, kurz ASC")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($kurz, $name, $stufe, $stamm, $punkte);
    while ($stmt->fetch()) {
      if($prev_punkte != $punkte) {
        $position++;
	    }
      if($stufenpunkte[$stufe] != $punkte) {
		    $stufenwertung[$stufe]++;
	    }
	    $stufenpunkte[$stufe] = $punkte;
      echo '<li><a href="#">'.$position.'. Platz ('.$stufenwertung[$stufe].'. der '.$Stufe[$stufe].') "'.$name.'" - '.$stamm.'<span class="ui-li-count">'.$punkte.'</span></a></li>';
      $prev_punkte = $punkte;
      $prev_stufe = $stufe;
    }
  }
  echo '</ul>';
}
else {
    echo 'Keine Berechtigung.';
}
