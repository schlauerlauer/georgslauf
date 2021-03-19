<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
session_start();
$user_check=$_SESSION['login_user'];
$row;
if ($stmt = $mysqli->prepare("SELECT username FROM login WHERE username = ?")) {
	$stmt->bind_param('s', $user_check);
	$stmt->execute();
	$stmt->store_result();
	$stmt->bind_result($row);
	while ($stmt->fetch()) {}
}
$login_session = $row;
if(!isset($row)){
	header('Location: /');
}
?>