<?php
include_once '../../includes/connect_gl.php';
include_once '../session/session.php';
include_once '../host/settings.php';

if($Abmeldung == true) {
if (isset ($login_session)) {
	if (isset ($_POST['pid'])) {
		if ($_POST['type'] == "p") $insert = "posten";
		else $insert = "gruppen";
		if ($stmt = $mysqli->prepare("DELETE FROM `".$insert."` WHERE `id` = ? AND `stamm` = ?")) {
			$stmt->bind_param('ss', $_POST['pid'], $login_session);
			$stmt->execute();
			echo "ok";
		} else echo "error";
	} else echo "error";
}
else echo "error";
}
else echo "nein";
?>
