<?php
include_once '../host/settings.php';
require('../session/session.php');
?>
<h1>Angemeldete Gruppen</h1>
<h2>Stamm <?php echo $login_session; ?></h2>
<div data-role="collapsible-set">
	<?php
	if($_SESSION['rolle'] >= 2) {
		if(isset($login_session)) {
			$user = $login_session;
			$gruppen = 0;
			if ($stmt = $mysqli->prepare("SELECT id, name, stufe, stamm, size, veggies FROM gruppen WHERE stamm = ? ORDER BY stufe, name asc")) {
				$stmt->bind_param('s', $user);
				$stmt->execute();
				$stmt->store_result();
				$stmt->bind_result($id, $name, $stufe, $stamm, $size, $veggies);
				while ($stmt->fetch()) {
					$gruppen++;
					echo '<div data-role="collapsible"><h3>'.$name.' - '.$Stufe[$stufe].'</h3><p>Stamm: '.$stamm.'<br><br>Größe: '.$size.' (Davon '.$veggies.' Vegetarier)</p>
					<a href="#" class="ui-btn ui-btn-inline ui-icon-delete ui-btn-icon-notext delete" type="g" name='.$name.' id='.$id.'></a>
					</div>';
				}
				if ($gruppen == 0) echo "Noch keine Gruppen angemeldet";
			}
		}
	} else echo "Keine Rechte";
	?>
</div>
<br>
<br>
<br>
<p align="center"><img src="../res/logo.png" style="max-width: 100%; max-height: 500px;"/></p>
