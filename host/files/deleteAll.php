<?php
include_once '../../../includes/connect_gl.php';
require('../../session/session.php');

if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo "<h2>We need more lemon pledge...</h2>"
  if ($stmt = $mysqli->prepare("DELETE FROM gruppen, posten")) {
    $stmt->execute();
    echo "Alle Gruppen und Posten gelöscht.";
  }
  else echo "Etwas ist schiefgelaufen";

} else echo "Keine Berechtigung";
?>
