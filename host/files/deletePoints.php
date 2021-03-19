<?php
include_once $__SERVER["DOCUMENT_ROOT"].'/includes/connect_gl.php';
require('../../session/session.php');

if(isset($login_session) && $_SESSION['rolle'] >= 4) {

  if ($stmt = $mysqli->prepare("DELETE FROM punkte")) {
    $stmt->execute();
    echo "Alle Bewertungen gelöscht.";
  }
  else echo "Etwas ist schiefgelaufen";

} else echo "Keine Berechtigung";
?>
