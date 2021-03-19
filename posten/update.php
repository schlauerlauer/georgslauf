<?php
require('../session/session.php');
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
include_once '../host/settings.php';
if(isset($login_session) && $_SESSION['rolle'] == 1) {
if($Punkte[0] == true) {
if(isset($_POST['p'])) {
  $user = $login_session;
  $pVal = $_POST['p'];
  if ($pVal > 100) $pVal = 100;
  else if ($pVal < 0) $pVal = 0;
  $vergleich = "< 2";
  if($user >= $PostenStufe) $vergleich = ">= 2";
  if ($stmt = $mysqli->prepare("INSERT INTO punkte (von, an, points) VALUES ((SELECT id FROM posten WHERE kurz = ?), IF((SELECT stufe FROM gruppen WHERE id = ?)".$vergleich.", ?, -1), ?)
  ON DUPLICATE KEY UPDATE points=VALUES(points), time=CURRENT_TIMESTAMP")) {
    $stmt->bind_param('siii', $user, $_POST['id'], $_POST['id'], $pVal);
    $stmt->execute();
    echo "ok";
  }
  else echo "error";
} else echo "error";
} else echo "error";
} else echo "error";
 ?>
