<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/includes/connect_gl.php';
include_once '../host/settings.php';
include_once '../session/session.php';
?>
<h1>Angemeldete Posten<h1>
<h2>Stamm <?php echo $login_session; ?></h2>
<div data-role="collapsible-set">
	<?php
		if($_SESSION['rolle'] >= 2) {
		if(isset($login_session)) {
			$posten = 0;
			if ($stmt = $mysqli->prepare("SELECT posten.id, kurz, password, color, name, kategorie, beschreibung, kontakt, anzahl, veggie, material, ort, sonstiges, stufe
				FROM posten, login
				WHERE stamm = ? AND kurz = username
				ORDER BY stufe DESC, name ASC")) {
				$stmt->bind_param('s', $login_session);
				$stmt->execute();
				$stmt->store_result();
				$stmt->bind_result($id, $kurz, $pw, $color, $name, $kat, $desc, $kontakt, $val, $veg, $mat, $ort, $sonst, $stufe);
				while ($stmt->fetch()) {
					$posten++;
					echo '<div data-role="collapsible"><h3>';
					if($ShowLogin == true) echo $kurz.' &nbsp;';
					echo $stufe.' - '.$Kat[$kat].' - '.$name.'</h3><p>';
					if($ShowLogin == true) echo '<h3 class="ui-bar ui-bar-a ui-corner-all" style="background-color:'.$Farben[$color].';" >Login '.$kurz.'<br />Passwort '.$pw.'</h3><br />';
					echo '"'.$desc.'"<br><br>Kontakt <strong>'.$kontakt.'</strong><br>Leiter <strong>'.$val.'</strong> (davon Veggie <strong>'.$veg.')</strong><br>';
					if (!empty($mat)) echo '<br>benötigtes Material <strong>'.$mat.'</strong>';
					if (!empty($ort)) echo '<br>bevozugter Ort <strong>'.$ort.'</strong>';
					if (!empty($sonst)) echo '<br>Sonstiges <strong>'.$sonst.'</strong>';
					echo '<br><a href="#" class="ui-btn ui-btn-inline ui-icon-delete ui-btn-icon-notext delete" type="p" name='.$name.' id='.$id.'></a></p></div>';
				}
				if ($posten == 0) echo "Noch keine Posten angemeldet";
			}
		}
	} else echo "Keine Rechte";
	?>
</div>
<br>
<br>
<br>
<p align="center"><img src="../res/logo.png" style="max-width: 100%; max-height: 500px;"/></p>
