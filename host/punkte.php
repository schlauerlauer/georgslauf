<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
include_once 'settings.php';
require('../session/session.php');

if(isset($login_session) && $_SESSION['rolle'] >= 3) {

  if(isset($_POST['punkte'])) {
    if ($stmt = $mysqli->prepare("UPDATE punkte SET points = ?, time = CURRENT_TIMESTAMP WHERE von = ? AND an = ?")) {
      $stmt->bind_param('sss', $_POST['punkte'], $_POST['von'], $_POST['an']);
      $stmt->execute();
      echo "ok";
    }
    else echo "error";
  }

}
?>
