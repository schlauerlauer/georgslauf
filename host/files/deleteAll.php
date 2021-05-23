<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
require('../../session/session.php');

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  if ($stmt = $mysqli->prepare("DELETE FROM gruppen, posten")) {
    $stmt->execute();
    echo "Alle Gruppen und Posten gelöscht";
  }
  else echo "Etwas ist schiefgelaufen";

} else echo "Keine Berechtigung";
?>
