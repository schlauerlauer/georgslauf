<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo "<h2>Eigene Query ausführen</h2>";

  $val = 0;
  $anzahl = 2;
  $zeile = "SELECT kurz";
  if($anzahl > 1) $zeile .= ", kontakt";
  $zeile .= " FROM posten";
  echo '<ol>';
  if ($stmt = $mysqli->prepare($zeile)) {
    $stmt->execute();
    $stmt->store_result();
    if($anzahl == 1) $stmt->bind_result($eins);
    if($anzahl == 2) $stmt->bind_result($eins, $zwei);
    //$stmt->bind_result($eins);
    while ($stmt->fetch()) {
      $val++;
      echo '<li>'.$eins.' '.$zwei.'</li>';
    }
  }
  echo '</ol>';
  echo "<br>Insgesamt ".$val.' Zeilen';
} else echo "Keine Berechtigung";

?>
