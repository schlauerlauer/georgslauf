<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/includes/connect_gl.php';
include_once '../host/settings.php';
include_once '../session/session.php';
include_once '../host/mail.php';
require './pGet.php';

if($Anmeldung == true) {
if (isset ($_POST['gorp'])) {
	if ($_POST['gorp'] == "g_save") {
		if ($stmt = $mysqli->prepare("INSERT INTO `gruppen` (`name`, `size`, `stufe`, `stamm`, `veggies`) VALUES (?,?,?,?,?)")) {
			$stmt->bind_param('sssss', $_POST['gname'], $_POST['gval'], $_POST['stid'], $login_session, $_POST['gveg']);
			$stmt->execute();
			echo 1;
			sendmail('Gruppenanmeldung von '.$login_session.' ('.$Stufe[$_POST['stid']].')','Neue Gruppe angemeldet von Stamm '.$login_session.' mit '.$_POST['gval'].' Kindern');
		} else echo 2;
	} else if ($_POST['gorp'] == "p_save") {
		if ($_POST['pid'] == "RoPo") {
			if ($rP[$_POST['kid']] >= $RKat[$_POST['kid']]) {
				echo 2;
				return;
			}
		} else {
			if ($wP[$_POST['kid']] >= $WKat[$_POST['kid']]) {
				echo 2;
				return;
			}
		}
		if ($stmt = $mysqli->prepare("INSERT INTO `posten` (`name`,`stamm`,`kategorie`,`beschreibung`,`kontakt`,`anzahl`,`veggie`,`material`,`ort`,`sonstiges`,`stufe`) VALUES (?,?,?,?,?,?,?,?,?,?,?)")) {
			$stmt->bind_param('sssssssssss', $_POST['pname'], $login_session, $_POST['kid'], $_POST['pdesc'], $_POST['pkont'], $_POST['pval'], $_POST['pveg'], $_POST['pmat'], $_POST['port'], $_POST['psonst'], $_POST['pid']);
			$stmt->execute();
			echo 1;
			sendmail('Postenanmeldung von '.$login_session.' ('.$_POST['pid'].')','Neuer Posten angemeldet von Stamm '.$login_session.', Kategorie '.$Kat[$_POST['kid']].', Name '.$_POST['pname']);
		} else echo "error";
	} else echo "error";
}
else echo "error";
}
else echo "Anmeldung nur noch über email möglich";
?>
