<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/master/includes/connect_gl.php';
require('../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {

  if(isset($_POST['punkte'])) {
    if ($stmt = $mysqli->prepare("DELETE from punkte where id = ?")) {
			$stmt->bind_param('s', $_POST['punkte']);
			$stmt->execute();
			echo "ok";
		} else echo "error";
  }

}
?>
