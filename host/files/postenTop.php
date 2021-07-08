<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/master/includes/connect_gl.php';
require('../../session/session.php');
include_once '../settings.php';

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo "<h2>Posten Sieger</h2>";
  $position = 0;
  echo '<ul data-role="listview" data-count-theme="b" data-inset="true">';
  if ($stmt = $mysqli->prepare("SELECT kurz, name, stufe, stamm, sum(points) as summe FROM posten, punkte WHERE posten.id = an GROUP BY an ORDER BY summe DESC")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($kurz, $name, $stufe, $stamm, $punkte);
    while ($stmt->fetch()) {
      if($prev_punkte != $punkte) $position++;
      echo '<li><a href="#">'.$position.'. Platz - '.$kurz.' - "'.$name.'" - '.$stamm.'<span class="ui-li-count">'.$punkte.'</span></a></li>';
      $prev_punkte = $punkte;
    }
  }
  echo '</ul>';
}
else {
    echo 'Keine Berechtigung.';
}