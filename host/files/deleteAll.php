<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
require('../../session/session.php');

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  if ($stmt = $mysqli->prepare("DELETE FROM gruppen")) {
    $stmt->execute();
    echo "Alle Gruppen gelöscht.";
  } else echo "Fehler beim Gruppen löschen.";
  if ($stmt = $mysqli->prepare("DELETE FROM posten")) {
    $stmt->execute();
    echo "Alle Posten gelöscht.";
  } else echo "Fehler beim Posten löschen.";
} else echo "Keine Berechtigung";
?>
