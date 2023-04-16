<?php
include_once '../host/settings.php';
require('../session/session.php');
include_once '../host/mail.php';
require './pGet.php';

if($Anmeldung == true) {
if (isset ($_POST['gorp'])) {
	if ($_POST['gorp'] == "g_save") {
		if ($stmt = $mysqli->prepare("INSERT INTO `gruppen` (`name`, `size`, `stufe`, `stamm`, `veggies`, `kontakt`) VALUES (?,?,?,?,?,?)")) {
			$stmt->bind_param('ssssss', $_POST['gname'], $_POST['gval'], $_POST['stid'], $login_session, $_POST['gveg'], $_POST['gnum']);
			$stmt->execute();
			echo 1;
			sendmail('Gruppenanmeldung von '.$login_session.' ('.$Stufe[$_POST['stid']].')','Neue Gruppe angemeldet von Stamm '.$login_session.' mit '.$_POST['gval'].' Kindern');
		} else echo 2;
	} else if ($_POST['gorp'] == "p_save") {
		if ($current_posten_pro_kategorie[$_POST['kid']] >= $max_posten_pro_kategorie[$_POST['kid']]) {
			echo 2;
			return;
		}
		echo 3;
		return;
		if ($stmt = $mysqli->prepare("INSERT INTO `posten` (`name`,`stamm`,`kategorie`,`beschreibung`,`kontakt`,`anzahl`,`veggie`,`material`,`ort`,`sonstiges`,`stufe`) VALUES (?,?,?,?,?,?,?,?,?,?,?)")) {
			$stmt->bind_param('sssssssssss', $_POST['pname'], $login_session, $_POST['kid'], $_POST['pdesc'], $_POST['pkont'], $_POST['pval'], $_POST['pveg'], $_POST['pmat'], $_POST['port'], $_POST['psonst'], $_POST['pid']);
			$stmt->execute();
			echo 1;
			sendmail('Postenanmeldung von '.$login_session.' Neuer Posten angemeldet von Stamm '.$login_session.', Kategorie '.$Kat[$_POST['kid']].', Name '.$_POST['pname']);
		} else echo "error";
	} else echo "error";
}
else echo "error";
}
else echo "Anmeldung nur noch über email möglich";
?>
