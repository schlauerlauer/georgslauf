<?php
include_once '../includes/connect_gl.php';
session_start();
$row = '';
if (isset($_POST['un'])) {
	if (empty($_POST['un']) || empty($_POST['pw'])) {
	} else {
		$username = $_POST['un'];
		$password = $_POST['pw'];
		if ($stmt = $mysqli->prepare("SELECT username, rolle from login WHERE password = ? AND username = ? LIMIT 1")) {
			$stmt->bind_param('ss', $password, $username);
			$stmt->execute();
			$stmt->store_result();
			$stmt->bind_result($row, $rolle);
			while ($stmt->fetch()) {}
			if ($_POST['un'] == $row) {
				$_SESSION['login_user'] = $username;
				$_SESSION['rolle'] = $rolle;
				echo $_SESSION['rolle'];
			}
			else {
			}
		}
	}
}
	//stripslashes?
	//mysqli_real_escape_string
?>
